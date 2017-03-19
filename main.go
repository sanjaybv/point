package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	ws "github.com/gorilla/websocket"
)

const defaultPort = "8080"

// upgrades HTTP connections to Websocket
var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	// for heroku support
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/index.html", index)
	http.HandleFunc("/ws", handleWs)

	log.Println("starting server on :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
	}
}

func handleWs(w http.ResponseWriter, r *http.Request) {

	// upgrade HTTP to Websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	NewClient(conn)
}

func index(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("public/index.html")
	if err != nil {
		log.Fatal(err)
	}

	if err := t.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}
