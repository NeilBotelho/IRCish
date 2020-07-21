package main

import(
	"math/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"
)
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