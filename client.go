package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
)

type Client struct {
	identity  string         //Needed?
	writeCh   *chan Msg      // send recieve message from broadcaster
	terminate *chan struct{} // terminate signal
	conn      *websocket.Conn
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade http to websocket
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Make write channel and terminate channel for client object
	writeCh := make(chan Msg)
	terminate := make(chan struct{})
	// Make client object
	client := Client{
		writeCh: &writeCh,
		terminate: &terminate,
		conn:      conn,
		identity:  strconv.Itoa(RandInt()), //Random Integer
	}

	// Register client writer in a go routine
	go clientWriter(&client)

	//Announce creation of client to broadcaster
	msg := Msg{
		OpCode:  &notify,
		Content: client.identity + " Just entered",
		Room:    defaultRoom,
		client:  &client,
	}
	messaging <- msg
	// Optionally set Opcode to join, but not necessary
	entering <- msg
	for {
		//Clear msg
		msg=Msg{}
		
		_, clientMsg, err := client.conn.ReadMessage()
		if err != nil {
			// Close connection in case of error
			leaving <- Msg{
				OpCode: &leaveAll,
				client: &client,
				Content:client.identity+" Just left",
			}
			*client.terminate <- struct{}{}
			close(*client.writeCh)
			close(*client.terminate)
			log.Println("Socket closed",err)
			return
		}
		// Set initial ReadDeadline
		conn.SetReadDeadline(time.Now().Add(time.Second * pingTimeout))
		err = json.Unmarshal(clientMsg, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling")
			continue
		}
		resolveRequest(&client,msg)
	}

}
func resolveRequest(client *Client, msg Msg) {
	msg.client=client
	switch *msg.OpCode {
		case communicate:
			log.Println("Communication")
			msg.From = &client.identity
			messaging <- msg
		case join:
			log.Println("Join")
			
			msg.OpCode = &notify
			msg.Content=client.identity+" Just joined"
			messaging <- msg
			msg.OpCode= &join
			entering <- msg
		case leave:
			leaving <- msg
			log.Println("Leaving")
		case identify:
			log.Println("Identify")
			if usernameValidate(msg.Content) {
				messaging <- Msg{
					OpCode:  &notifyAll,
					Content: client.identity + " --> " + msg.Content,
					client: client,
				}
				client.identity = msg.Content
			}
		case ping:
			log.Println("Ping")
			return
	}
}

func clientWriter(cli *Client) {
	pinger:=time.Tick(+(time.Second*pingTimeout-4))
	for {
		select {
		case <-*cli.terminate:
			close(pinger)
			return
		case msg := <-*cli.writeCh:
			out, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			cli.conn.WriteMessage(websocket.TextMessage, out)
		case <-pinger:
			cli.conn.WriteMessage(websocket.TextMessage,[]byte(`{"opcode":4}`))
		}
	}
}

func RandInt() int {
	return rand.Intn(89999) + 10000
}


func usernameValidate(username string) bool {
	var pass bool
	var err error
	if pass, err = regexp.MatchString(`\s`, username); pass && err == nil {
		return false
	}
	pass, err = regexp.MatchString(`[\w]{2,10}`, username)
	if pass, err = regexp.MatchString(`[\W]{2,10}`, username); pass && err == nil {
		return false
	}
	return true
}
