// TODO
// 1. Simplify the http adding the handler part
// 		and check if it all works out
// 2. Write a Json type
// 3. Do something about the Client type
package comm

import (
	ws "golang.org/x/net/websocket"

	"container/list"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type CommService struct {
	funcHandler func(*Client)
	conns       *list.List
}

type Client struct {
	conn *ws.Conn
	comm *CommService
}

func NewCommService(url, port string, funcHandler func(*Client)) {

	cs := CommService{funcHandler, list.New()}

	http.HandleFunc("/"+url,
		func(w http.ResponseWriter, req *http.Request) {
			s := ws.Server{
				Handler: ws.Handler(
					cs.newClientsHandler)}
			s.ServeHTTP(w, req)
		})
}

func (cs *CommService) newClientsHandler(conn *ws.Conn) {

	cl := &Client{conn, cs}

	connElem := cs.conns.PushBack(cl)
	defer cs.conns.Remove(connElem)

	cs.funcHandler(cl)
}

func (c *CommService) BroadcastJson(msg map[string]interface{}) {

	for elem := c.conns.Front(); elem != nil; elem = elem.Next() {
		go func(elem *list.Element) {
			if err := elem.Value.(*Client).SendJson(msg); err != nil {
				log.Println("broadcast:", err)
				c.conns.Remove(elem)
			}
		}(elem)
	}
}

func Connect(host, port, serviceName string) (*Client, error) {

	origin := "http://" + host + "/"
	url := "ws://" + host + ":" + port + "/" + serviceName
	if conn, err := ws.Dial(url, "", origin); err != nil {
		return nil, errors.New(fmt.Sprintln(
			"comm: websocket connect error:", err))
	} else {
		return &Client{conn, nil}, nil
	}
}

func (cl *Client) RecvJson() (map[string]interface{}, error) {

	data := make(map[string]interface{})
	err := ws.JSON.Receive(cl.conn, &data)
	if err != nil {
		e := errors.New(fmt.Sprintln("comm: websocket receive error:", err))
		return nil, e
	}

	return data, nil
}

func (cl *Client) SendJson(msg map[string]interface{}) error {

	err := ws.JSON.Send(cl.conn, msg)
	if err != nil {
		e := errors.New(fmt.Sprintln("comm: websocket send error:", err))
		return e
	}

	return nil
}

func (cl *Client) GetCommService() *CommService {
	return cl.comm
}
