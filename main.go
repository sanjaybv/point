package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"github.com/streamrail/concurrent-map"
)

const defaultPort = "8080"

// upgrader is used to upgrade HTTP connections to Websocket
var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var broadcast chan []byte
var connPool cmap.ConcurrentMap

func main() {

	// for heroku support
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/index.html", index)
	http.HandleFunc("/ws", handleWs)

	connPool = cmap.New()
	broadcast = make(chan []byte)
	go broadcastWorker(broadcast)

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

	connID := uuid.NewV4().String()
	connPool.Set(connID, conn)
	defer connPool.Remove(connID)

	log.Println("client", "new", connID, r.RemoteAddr)

	// PingWorker send pings to keep the connection alive
	go func() {
		for {
			if _, ok := connPool.Get(connID); ok {
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

	// read messages and send it to broadcast channel
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

// BroadcastWorker sends the message to all connections in connPool
// in its own goroutine.
func broadcastWorker(broadcast chan []byte) {
	log.Println("broadcastWorker running")

	for {
		// get the data
		data := <-broadcast

		// broadcast the data
		for connTuple := range connPool.IterBuffered() {
			go func(conn *ws.Conn, data []byte) {
				if err := conn.WriteMessage(ws.TextMessage, data); err != nil {
					log.Println(err)
				}
			}(connTuple.Val.(*ws.Conn), data)
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
