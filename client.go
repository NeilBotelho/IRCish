package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
)

const (
	// Channel Buffer Size constants (never changed)
	clientMsgBuff uint8 = 1
	chanBuff      uint8 = 1
	// Ping timeout
	pingTimeout = 15 //seconds
)

var (
	// operations
	communicate uint8  = 0
	join        uint8  = 1
	leave       uint8  = 2
	identify    uint8  = 3
	ping        uint8  = 4
	leaveAll    uint8  = 5
	defaultRoom string = "general"

	//Communication Channel
	entering  = make(chan Msg,chanBuff)
	leaving   = make(chan Msg,chanBuff)
	messaging = make(chan Msg,chanBuff)
)

type Client struct {
	identity  string         //Needed?
	writeCh   *chan Msg      // send recieve message from broadcaster
	terminate *chan struct{} // terminate signal
	conn      *websocket.Conn
}

type Msg struct {
	OpCode  *uint8 `json:"opcode"`
	Content string `json:"content,omitempty"`
	Room    string `json:"room,omitempty"`
	client  *Client
	From    *string `json:"from,omitempty"`
}

type Room map[*Client]bool

func broadcaster() {
	RoomList := make(map[string]Room)
	for {
		select {
		case msg := <-messaging:
			if RoomList[msg.Room] != nil {
				fmt.Println("Currentl not nil")
				for cli, _ := range RoomList[msg.Room] {

					fmt.Printf("Sending message %s on room %s to client %s\n", msg.Content, msg.Room, msg.client.identity)
					*cli.writeCh <- msg
				}
			}
		case msg := <-entering:
			if RoomList[msg.Room] == nil {
				RoomList[msg.Room] = Room{}
			}
			RoomList[msg.Room][msg.client] = true

		case msg := <-leaving:
			fmt.Printf("leaving with code %d\n", *msg.OpCode)
			if *msg.OpCode == leaveAll {
				for roomName, _ := range RoomList {
					fmt.Printf("client %s Leaving room %s", msg.client.identity, roomName)
					delete(RoomList[roomName], msg.client)
					if len(RoomList[roomName]) == 0 {
						delete(RoomList, roomName)
					}
				}
			} else {
				delete(RoomList[msg.Room], msg.client)
				if len(RoomList[msg.Room]) == 0 {
					delete(RoomList, msg.Room)
				}
			}
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
	// Make write and terminate channel
	writeCh := make(chan Msg)
	terminate := make(chan struct{})

	// Make client object
	client := Client{writeCh: &writeCh,
		terminate: &terminate,
		conn:      conn,
		identity:  strconv.Itoa(RandInt()), //Random Integer
	}

	//  Create Msg object to be used by this thread
	var msg Msg

	// Register client writer in a go routine
	go clientWriter(&client)

	//Announce creation of client to broadcaster
	msg = Msg{
		OpCode:  &communicate,
		Content: client.identity + " Just entered",
		Room:    defaultRoom,
		client:  &client,
		From:    &client.identity,
	}
	messaging <- msg
	fmt.Println("Entering")
	entering <- msg

	//Set initial ReadDeadline
	// conn.SetReadDeadline(time.Now().Add(time.Second * pingTimeout))

	for {
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			// Close connection in case of error
			leaving <- Msg{OpCode: &leaveAll, client: &client}
			*client.terminate <- struct{}{}
			close(writeCh)
			close(terminate)

			log.Println("Socket closed")
			return
		}
		// conn.SetReadDeadline(time.Now().Add(time.Second * pingTimeout))
		err = json.Unmarshal(clientMsg, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling")
			continue
		}

		switch msg.OpCode {
		case &communicate:
			msg.From = &client.identity
			messaging <- msg
			fmt.Println("Communication")
		case &join:
			entering <- msg
			msg.OpCode = &communicate
			messaging <- msg
		case &leave:
			leaving <- msg
			fmt.Println("leaving")
			// case &identify:
			// if nameCheck(msg.)
			// client.identity=msg.content

		}
	}

}

func clientWriter(cli *Client) {
	for {
		select {
		case <-*cli.terminate:
			return
		case msg := <-*cli.writeCh:
			out, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			cli.conn.WriteMessage(websocket.TextMessage, out)
		}
	}
}

func RandInt() int {
	return rand.Intn(89999) + 10000
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
