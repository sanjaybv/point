package main

import (
	"log"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/streamrail/concurrent-map"
)

var connPool cmap.ConcurrentMap

// Client represents a client
type Client struct {
	conn   *ws.Conn
	connID string

	writeChan chan []byte
}

// NewClient creates a new client object with the given websocket connection
func NewClient(c *ws.Conn) *Client {
	client := &Client{
		conn:      c,
		connID:    c.RemoteAddr().String(),
		writeChan: make(chan []byte),
	}
	go client.pingWorker()
	go client.writeWorker()
	go client.broadcastWorker()

	log.Println("client", "new", client.connID)
	connPool.Set(client.connID, client)

	return client
}

// Close closes the underlying websocket connection
func (c *Client) Close(from string) {

	connPool.Remove(c.connID)
	c.conn.Close()
	log.Println("client", "exit-from", from, c.connID)
}

func (c *Client) pingWorker() {

	for {
		if err := c.conn.WriteControl(
			ws.PingMessage, nil, time.Now().Add(5*time.Second)); err != nil {
			log.Println(err)
			break
		}
		time.Sleep(5 * time.Second)
	}

	c.Close("ping")
}

func (c *Client) writeWorker() {
	for {
		data := <-c.writeChan
		if err := c.conn.WriteMessage(ws.TextMessage, data); err != nil {
			log.Println(err)
			break
		}
	}

	c.Close("write")
}

func (c *Client) broadcastWorker() {
	for {
		// read
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		// broadcast
		for connTuple := range connPool.IterBuffered() {
			connTuple.Val.(*Client).writeChan <- data
		}
	}

	c.Close("broadcast")
}
