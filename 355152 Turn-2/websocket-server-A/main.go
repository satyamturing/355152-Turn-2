package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketServerURL = "ws://localhost:8080/ws"  
	restAPIURL         = "http://localhost:8080/api"  
)

type Message struct {
	Text string `json:"text"`
}

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

	go receiveMessages(wsConn)
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

func receiveMessages(wsConn *websocket.Conn) {
	for {
		_, messageBytes, err := wsConn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading message from WebSocket server:", err)
		}

		var message Message
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}

		fmt.Println("Received message from WebSocket server:", message.Text)

		// Make an HTTP request to the REST API
		response, err := http.Get(restAPIURL)
		if err != nil {
			log.Fatal("Error making HTTP request:", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Fatalf("HTTP request failed with status: %s", response.Status)
		}

		fmt.Println("REST API response status:", response.Status)

		// Reading the response body
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Error reading response body:", err)
		}

		fmt.Println("REST API response body:", string(body))
	}
}