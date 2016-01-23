package cserver

import (
	"sanjay/point/comm"

	"log"
	"time"
)

const (
	PORT = "8080"
	//TIME_FORMAT = "Jan 2, 2006 at 3:04pm (MST)"
	TIME_FORMAT = "3:04:05pm"
)

func InitChatService() {

	comm.NewCommService("chat", ":"+PORT, newClient)

	log.Println("cserver: started...")
}

func newClient(c *comm.Client) {

	log.Println("server: new client")

	for {
		lastTime := time.Now()

		// Receive the message
		msg, err := c.RecvJson()
		if err != nil {
			log.Println(err)
			return
		}

		// Process the message

		if time.Now().Sub(lastTime) < time.Duration(500*time.Millisecond) {
			continue
		}

		msg["time"] = time.Now().Format(TIME_FORMAT)

		// Broadcast the message
		c.GetCommService().BroadcastJson(msg)
	}

}
