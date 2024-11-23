package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketServerURL = "ws://localhost:8080/ws" // Replace this with the actual WebSocket server URL
	restAPIURL         = "http://localhost:8080/api" // Replace this with the actual REST API endpoint URL
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// Connect to the WebSocket server
	wsConn, _, err := wsUpgrader.Upgrade(nil, &http.Request{URL: strings.NewReader(websocketServerURL)}, nil)
	if err != nil {
		log.Fatal("Error upgrading to WebSocket:", err)
	}
	defer wsConn.Close()

	// Send a message to the WebSocket server
	err = wsConn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket Server!"))
	if err != nil {
		log.Fatal("Error sending message to WebSocket server:", err)
	}

	fmt.Println("Sent message to WebSocket server.")

	// Read message from WebSocket server
	_, message, err := wsConn.ReadMessage()
	if err != nil {
		log.Fatal("Error reading message from WebSocket server:", err)
	}

	fmt.Println("Received message from WebSocket server:", string(message))

	// Make a request to the REST API endpoint
	response, err := http.Get(restAPIURL)
	if err != nil {
		log.Fatal("Error making HTTP request:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed with status: %s", response.Status)
	}

	fmt.Println("Made a successful HTTP request to the REST API.")

	// Reading the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	fmt.Println("Response from REST API:", string(body))
}