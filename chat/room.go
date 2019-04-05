package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"p-bitbucket.imovetv.com/heracles/trace"
	"strconv"
)

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients
	forward chan []byte
	// join is a channel for clients wishing to join the room
	join chan *client
	// leave is a channel for clients wishing to leave the room
	leave chan *client
	// clients holds all current clients in this room
	clients map[*client]bool
	// tracer will receive trace information of activity in room
	tracer trace.Tracer
}

// newRoom makes a new room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}
var clientNumber = 0

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	clientNumber++
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	for name, values := range req.Header {
		// Loop over all values in header
		for _, value := range values {
			r.tracer.Trace(name, " -- ", value)
		}
	}
	client := &client{
		name:   strconv.Itoa(clientNumber),
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
