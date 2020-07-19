package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var PORT string = "8000"

func main() {
	if len(os.Args) > 1 {
		PORT = os.Args[1]
	}

	r := mux.NewRouter() //.StrictSlash(false)
	srv := &http.Server{
		// Wrap mux in a timeout Handler
		Handler: http.TimeoutHandler(r, 10*time.Second, "Timeout!\n"),
		Addr:    "127.0.0.1:" + PORT,
		// Enforce timeouts for server!
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	r.HandleFunc("/test", HomeHandler2)
	r.HandleFunc("/{hello}", HomeHandler)
	// r.HandleFunc("/products", ProductsHandler)
	// r.HandleFunc("/articles", ArticlesHandler)
	log.Print("Server running on " + PORT)
	log.Fatal(srv.ListenAndServe())

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "Hello at %q", html.EscapeString(r.URL.Path))
}
func HomeHandler2(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprintf(w, "test at %q", html.EscapeString(r.URL.Path))
}
