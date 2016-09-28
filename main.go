package main

import (
	ws "github.com/gorilla/websocket"
	"github.com/satori/go.uuid"

	//"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var broadcast chan []byte
var connPool map[string]*ws.Conn

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// increase security on this
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/index.html", index)
	http.HandleFunc("/ws", handleWs)

	connPool = make(map[string]*ws.Conn)
	broadcast = make(chan []byte)
	go broadcastWorker(broadcast)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
	}
}

func handleWs(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connID := uuid.NewV4().String()
	connPool[connID] = conn
	defer delete(connPool, connID)

	log.Println("client", "new", connID, r.RemoteAddr)

	// send pings to keep the connection alive
	go func() {
		for {
			if _, ok := connPool[connID]; ok {
				if err := conn.WriteControl(
					ws.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
					log.Println(err)
				}
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		broadcast <- data
	}

	conn.Close()
	log.Println("client", "exit", connID)
}

func broadcastWorker(broadcast chan []byte) {
	log.Println("broadcastWorker running")

	for {
		data := <-broadcast
		//log.Println(string(data))
		for _, conn := range connPool {
			if err := conn.WriteMessage(ws.TextMessage, data); err != nil {
				log.Println(err)
			}
		}
	}
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
