package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
)

// Represents a client
// contains all information needed to 
// send and recieve messages to/from client
type Client struct {
	identity  string       	  
	writeCh   *chan Msg        // send recieve message from broadcaster
	terminate *chan struct{}   // used to send and recieve terminate signal
	conn      *websocket.Conn  // websocket connection with client
}


func clientCreator(w http.ResponseWriter, r *http.Request) {
/*Create Client object and announce client entering.
Runs clientHandler in a goroutine and exits*/

	// Upgrade http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Make write channel and terminate channel for Client object
	writeCh := make(chan Msg)
	terminate := make(chan struct{})
	
	// Make Client object
	client := Client{
		writeCh: &writeCh,
		terminate: &terminate,
		conn:      conn,
		identity:  strconv.Itoa(RandInt()), //Random Integer
	}


	//Announce creation of client to broadcaster
	msg := Msg{
		OpCode:  &notify,
		Content: client.identity + " Just entered",
		Room:    defaultRoom,
		client:  &client,
	}
	messaging <- msg

	// Add client to defaultRoom
	entering <- msg
	go clientHandler(&client)
}

func clientHandler(client *Client){
/*Listens for incoming messages from clients and handles them*/

	// Register client writer in a go routine
	go clientWriter(client)

	var msg Msg
	for {
		//Clear previous contents of msg
		msg=Msg{}
		
		// Wait for client message(sleeps if no message) 
		// ReadDeadline set in resolveRequest
		_, clientMsg, err := client.conn.ReadMessage()
		
		if err != nil {
			// Close connection in case read of error
			closeClient(client)
			log.Println("Socket closed",err)
			return
		}

		// Unpack JSON message from client into msg struct 
		err = json.Unmarshal(clientMsg, &msg)
		if err != nil {
			log.Println("Error unmarshalling")
			continue
		}

		// Reolve client request depending on opcode
		resolveRequest(client,msg)
	}

}
func resolveRequest(client *Client, msg Msg) {
/* Performs required operation depending on what opCode is sent by client*/
	
	// Set initial ReadDeadline
	// this closes the connection if there is no message from 
	// the client within pingTimeout seconds  
	client.conn.SetReadDeadline(time.Now().Add(time.Second * pingTimeout))
	
	// Message will not contain information about client so set it,
	// incase the operation requires it
	msg.client=client
	
	switch *msg.OpCode {
		case communicate:
			log.Println("Communication")
			// Set From field and send on messaging channel
			msg.From = &client.identity
			messaging <- msg
		case join:
			log.Println("Join")
			// Set opCode to notify and send on messaging channel
			msg.OpCode = &notify
			msg.Content=client.identity+" Just joined"
			messaging <- msg
			// Also add client to specified room by sending message on the entering channel
			msg.OpCode= &join
			entering <- msg
		case leave:
			// Send msg on leaving channel
			leaving <- msg
			log.Println("Leaving")
		case identify:
			log.Println("Identify")
			// If specified username is valid, 
			// change user identity and 
			// notify users in rooms current user is joined
			if usernameValidate(msg.Content) {
				messaging <- Msg{
					OpCode:  &notifyAll,
					Content: client.identity + " --> " + msg.Content,
					client: client,
				}
				client.identity = msg.Content
			}
		case ping:
			// This opcode is used to ensure client is still connected
			// and to prevent ReadDeadline from closing active clients
			log.Println("Ping")
			return
	}
}


func clientWriter(cli *Client) {
/*Recieves messages from cli.writeCh and send to client.  
Pings client periodically to prevent ReadDeadlline from closing active clients*/

	// Create a recieve channel that recieves a value every pingTimeout/2 seconds
	pinger:=time.NewTicker(+(time.Second*pingTimeout/2))
	for {
		select {
		case <-*cli.terminate:
			// If terminate signal recieved, close pinger and exit
			pinger.Stop()
			return
		case msg := <-*cli.writeCh:
			// If message is recieved, convert it to JSON and send to client 
			out, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			cli.conn.WriteMessage(websocket.TextMessage, out)
		case <-pinger.C:
			// If value is recieved from pinger, send a ping message to client
			cli.conn.WriteMessage(websocket.TextMessage,[]byte(`{"opcode":4}`))
		}
	}
}


func closeClient(client *Client){
/* Annouces client departure in all rooms client joined and 
closes any open channels or goroutines*/

	leaving <- Msg{
		OpCode: &leaveAll,
		client: client,
		Content:client.identity+" Just left",
	}

	//Send terminate signal(closes client writer)
	*client.terminate <- struct{}{} 
	
	// Close any open channels(to prevent memory leak)
	close(*client.writeCh)
	close(*client.terminate)
}


func RandInt() int {
/*Generate a random 5 digit number*/

	return rand.Intn(89999) + 10000
}


func usernameValidate(username string) bool {
/*Check if username is valid*/

	var pass bool
	var err error
	//Check if username contains whitespace
	if pass, err = regexp.MatchString(`\s`, username); pass && err == nil {
		return false
	}
	//Check that username has only alphanumeric characters and is between 2 and 10 chars 
	if pass, err = regexp.MatchString(`[\W]{2,10}`, username); pass && err == nil {
		return false
	}
	return true
}
