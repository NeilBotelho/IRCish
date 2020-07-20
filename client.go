package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	addr string
	ch chan []byte //send recieve message from broadcaster
	conn *websocket.Conn
	exit chan struct{} // terminate signal
}

var (
	entering=make(chan Client)
	leaving=make(chan Client)
	messages=make(chan []byte)
)

func broadcaster(){
	log.Println("broadcaster")
	clients:=make(map[Client]bool)
	for{
		select{
		case msg:= <-messages:
			for cli :=range clients{
				cli.ch<-msg
			}
		case cli := <-entering:
			clients[cli]=true;

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
	// Make communication and exit channel
	ch := make(chan []byte)
	exit := make(chan struct{})

	// Make client object
	client := Client{addr,ch, conn,exit}
	// Register client writer in a go routine
	go clientWriter(client)
	//Announce creation of client to broadcaster 
	entering<-client
	messages<- []byte(client.addr +" Just entered")
	for {
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Socket closed on client side")
			exit<-struct{}{}
			leaving<-client
			close(exit)
			close(ch)
			return
		}
		clientMsg=[]byte(addr+": "+string(clientMsg))
		log.Println(string(clientMsg))
		messages<-clientMsg
		log.Println("Message from " + r.RemoteAddr + ": " + string(clientMsg))
	}

}

func clientWriter(cli Client){
	for{
		select{
		case <-cli.exit:
			return
		case msg := <-cli.ch:
			cli.conn.WriteMessage(websocket.TextMessage,msg)
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
