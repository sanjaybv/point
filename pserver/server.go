package pserver

import (
	"github.com/sanjaybv/comm"

	"log"
	"sync/atomic"
	"time"
)

const (
	PORT = "8082"
)

var (
	msgCount int32
)

func InitPointService() {

	comm.NewCommService(":"+PORT, newClient)

	log.Println("pserver: started on 8082")

	go func() {
		intTime := time.Tick(time.Second)
		for {
			select {
			case <-intTime:
				cnt := atomic.LoadInt32(&msgCount)
				if cnt != 0 {
					log.Println("msgCount/second", cnt)
					atomic.StoreInt32(&msgCount, 0)
				}
			}
		}
	}()
}

func newClient(c *comm.Client) {

	log.Println("server: new client")

	for {
		// Receive the message
		msg, err := c.RecvJson()
		if err != nil {
			log.Println(err)
			return
		}

		// Process the message
		atomic.AddInt32(&msgCount, 1)

		// Broadcast the message
		c.GetCommService().BroadcastJson(msg)
	}
}
