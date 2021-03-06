package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

var PORT string = "8000"

// Creator the upgrader(upgrades http connection to websocket connection)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connection from other websites
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func main() {
	if len(os.Args) > 1 {
		PORT = os.Args[1]
	}
	r := mux.NewRouter() // maybe set .StrictSlash(false)?
	
	// Heroku Doesn't support the following hence it won't be used.
	// But it will be kept commented out as it is beneficial to have timeouts when they are supported.
	/*srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:" + PORT,
		// Enforce timeouts for server
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}*/

	// file server for static assets
	fs := http.FileServer(http.Dir("./static/")) 
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
	// Routes 
	r.HandleFunc("/ircish", clientCreator)
	r.HandleFunc("/", HomeHandler)

	go broadcaster()

	// Start server
	log.Print("Server running on " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT,r))

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	/*  Serves index.html */
	p,err:=ioutil.ReadFile("./public/index.html")
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", string(p))
}
