package main

import (
	"bytes"
	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {
	name string
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte
	// room is the room this client is chatting in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.send {
		var buffer bytes.Buffer
		buffer.WriteString(c.name)
		name := append([]byte("Client: "), buffer.Bytes()...)
		name = append(name, []byte(" - ")...)

		err := c.socket.WriteMessage(websocket.TextMessage, append(name, msg...))
		if err != nil {
			return
		}

	}
}
