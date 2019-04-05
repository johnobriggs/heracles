package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const (
	// NewRelicTxnKey is the key used to retrieve the NewRelic Transaction from the context
	NewRelicTxnKey = "NewRelicTxnKey"
)

// NewRelicMonitoring is a middleware that starts a newrelic transaction, stores it in the context, then calls the next handler
func NewRelicMonitoring(app newrelic.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		txn := app.StartTransaction(ctx.Request.URL.Path, ctx.Writer, ctx.Request)
		defer txn.End()
		ctx.Set(NewRelicTxnKey, txn)
		ctx.Next()
	}
}

func main() {

	engine := gin.Default()
	engine.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	engine.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello Gin Framework!"})
	})

	cfg := newrelic.NewConfig(os.Getenv("APP_NAME"), os.Getenv("NEW_RELIC_API_KEY"))
	_, err := newrelic.NewApplication(cfg)
	if err != nil {
		log.Printf("failed to make new_relic app: %v", err)
	} else {
		//engine.Use(adapters.NewRelicMonitoring(app))
	}

	//select
	//case value <- timer.c
	engine.Run(port())
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
}

// Get the port to listen on
func port() string {
	port, found := os.LookupEnv("PORT")
	if !found || len(port) == 0 {
		port = "8888"
	}
	return ":" + port
}

func echo(w http.ResponseWriter, r *http.Request) {

	message := r.URL.Query()["message"][0]
	w.Header().Add("Content-Type", "text-plain")
	fmt.Fprintf(w, message)

}

// Log the env variables required for a reverse proxy
func logSetup() {
	a_condtion_url := os.Getenv("A_CONDITION_URL")
	b_condtion_url := os.Getenv("B_CONDITION_URL")
	default_condtion_url := os.Getenv("DEFAULT_CONDITION_URL")

	log.Printf("Server will run on: %s\n", port())
	log.Printf("Redirecting to A url: %s\n", a_condtion_url)
	log.Printf("Redirecting to B url: %s\n", b_condtion_url)
	log.Printf("Redirecting to Default url: %s\n", default_condtion_url)
}

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) *json.Decoder {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		panic(err)
	}

	// Because go lang is a pain in the ass if you read the body then any subsequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

// Parse the requests body
func parseRequestBody(request *http.Request) requestPayloadStruct {
	decoder := requestBodyDecoder(request)

	var requestPayload requestPayloadStruct
	err := decoder.Decode(&requestPayload)

	if err != nil {
		panic(err)
	}

	return requestPayload
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}
