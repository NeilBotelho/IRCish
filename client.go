package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	addr string
	conn *websocket.Conn
}

// const(
// entering=make(chan string)

func wsHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Registered socket from " + r.RemoteAddr)

	// Read a message to this client
	// Blocks till message is recieved
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Socket closed on client side")
		return
	}
	log.Println("Message from " + r.RemoteAddr + ": " + string(message))

	// Send a message to client every second
	ticker := time.Tick(1 * time.Second)
	conn.WriteMessage(websocket.TextMessage, []byte("Hello socket"))
	messageNo := 0
	for {
		select {
		case <-ticker:
			messageNo++
			conn.WriteMessage(websocket.TextMessage, []byte("Message number "+strconv.Itoa(messageNo)))
		}
	}

}

func ws2Handler(w http.ResponseWriter, r *http.Request) {
	// Upgrade http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Get Address
	var addr string = r.Header.Get("X-FORWARDED-FOR")
	if addr == "" {
		addr = r.RemoteAddr
	}
	conn.WriteMessage(websocket.TextMessage, []byte("hi there client"))
	log.Println("Said hello to client")

	// Make client object
	client := Client{addr, conn}
	fmt.Println(client)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Socket closed on client side")
			return
		}
		log.Println("Message from " + r.RemoteAddr + ": " + string(message))
	}

}
