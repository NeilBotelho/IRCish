package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

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
	// Make write and terminate channel
	writeCh := make(chan Msg)
	terminate := make(chan struct{})

	// Make client object
	client := Client{writeCh: &writeCh,
		terminate: &terminate,
		conn:      conn,
		identity:  strconv.Itoa(RandInt()), //Random Integer
	}

	// Register client writer in a go routine
	go clientWriter(&client)

	//Announce creation of client to broadcaster
	msg := Msg{
		OpCode:  &communicate,
		Content: client.identity + " Just entered",
		Room:    defaultRoom,
		client:  &client,
		From:    &client.identity,
	}
	messaging <- msg
	fmt.Println("Entering")
	entering <- msg

	recieveAndResolve(&client)

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

func recieveAndResolve(client *Client) {

	//Set initial ReadDeadline
	// conn.SetReadDeadline(time.Now().Add(time.Second * pingTimeout))

	// Msg Var to be used in for loop
	var msg Msg

	for {
		_, clientMsg, err := client.conn.ReadMessage()
		if err != nil {
			// Close connection in case of error
			leaving <- Msg{OpCode: &leaveAll, client: client}
			*client.terminate <- struct{}{}
			close(*client.writeCh)
			close(*client.terminate)
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
		case &identify:
			if usernameValidate(msg.Content) {
				messaging <- Msg{
					OpCode:  &notify,
					Content: client.identity + " --> " + msg.Content,
				}
				client.identity = msg.Content
				*client.writeCh <- Msg{
					OpCode: &identify,
					Content: msg.Content
				}
			}
		}
	}
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
