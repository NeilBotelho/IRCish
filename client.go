package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// operations
	communicate   uint8 = 0
	join		  uint8 = 1
	leave         uint8 = 2
	identify      uint8 = 3
	ping          uint8 = 4

	// Channel Buffer Size constants
	clientMsgBuff uint8 = 1
	chanBuff      uint8 = 1

	// Ping timeout
	pingTimeout   uint8 = 15 //seconds	  
)

var (
	entering = make(chan Client,)
	leaving  = make(chan Client)
	messages = make(chan Msg)
)

type Client struct{ 
	identifier string //Needed?
	writeCh chan Msg // send recieve message from broadcaster
	terminate chan struct{} // terminate signal
	conn *websocket.Conn
}

type Msg struct{ 
	OpCode *uint8 `json:"opcode"`
	Content string `json:"content,omitempty"`
	Room string `json:"room,omitempty"`
}

type Room map[*Client]bool


func broadcaster() {
	var RoomList:=map[string]Room
	clients := make(map[Client]bool)
	for {
		select {
		case msg := <-messages:
			if *msg.OpCode == communicate {
				for cli := range clients {
					cli.ch <- msg
				}
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
		}
	}

}
func wsHandler(w http.ResponseWriter, r *http.Request) {
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
	addr += ":"
	// Make communication and exit channel
	ch := make(chan Msg)
	exit := make(chan struct{})

	// Make client object
	client := Client{addr, ch, conn, exit}

	// Register client writer in a go routine
	go clientWriter(&client)

	//Announce creation of client to broadcaster
	entering <- client
	messages <- Msg{&announce, client.addr + " Just entered"}

	var msg Msg
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			leaving <- client
			exit <- struct{}{}
			close(ch)
			close(exit)
			log.Println("Socket closed on client side")
			return
		}
		err = json.Unmarshal(clientMsg, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling")
		}
		messages <- msg
	}

}

func clientWriter(cli *Client) {
	for {
		select {
		case <-cli.exit:
			return
		case msg := <-cli.ch:
			cli.conn.WriteMessage(websocket.TextMessage, []byte(msg.Content))
		}
	}
}

func TestingHandler(w http.ResponseWriter, r *http.Request) {

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
