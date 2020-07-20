package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

var PORT string = "8000"
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	if len(os.Args) > 1 {
		PORT = os.Args[1]
	}

	r := mux.NewRouter() // maybe set .StrictSlash(false)?
	srv := &http.Server{
		// Wrap mux in a timeout Handler
		Handler: r,
		Addr:    "127.0.0.1:" + PORT,
		// Enforce timeouts for server
		// Will set longer when sockets implemented
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	r.HandleFunc("/ws", wsHandler)
	r.HandleFunc("/test", HomeHandler)
	log.Print("Server running on " + PORT)
	log.Fatal(srv.ListenAndServe())

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello at %q", html.EscapeString(r.URL.Path))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Registered socket from " + r.RemoteAddr)
	// time.Sleep(5 * time.Second)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Socket closed on client side")
			return
		}
		log.Println("Message from " + r.RemoteAddr + ": " + string(message))
		conn.WriteMessage(websocket.TextMessage, []byte("Hello socket"))
	}
	log.Println(conn)
}
