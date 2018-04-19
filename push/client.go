package push

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

type Client struct {
	Id          string
	Service     string
	ws          *websocket.Conn
	addClinetCh chan *Client
	rmClientCh  chan *Client
}

type Regist struct {
	Action  string `json:"action"`
	UserId  string `json:"userId"`
	Service string `json:"service"`
}

func NewClient(ws *websocket.Conn, add chan *Client, rm chan *Client) *Client {
	return &Client{
		ws:          ws,
		addClinetCh: add,
		rmClientCh:  rm,
	}
}

func (client *Client) Start() {
	for {
		var message string
		err := websocket.Message.Receive(client.ws, &message)
		if err != nil {
			client.rmClientCh <- client
			return
		} else {
			// json action:register, service:collie/shuffle, user:***
			// user regist
			fmt.Println(message)

			var a Regist
			err := json.Unmarshal([]byte(message), &a)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			if a.Action == "register" {
				client.Id = a.UserId
				client.Service = a.Service
				client.addClinetCh <- client
			}

		}
	}
}

/** notification */
func (client *Client) Send(message string) {
	err := websocket.Message.Send(client.ws, message)
	if err != nil {
		fmt.Println(message)
	}
}

func (client *Client) Close() {
	fmt.Println("client close")
	client.ws.Close()
}
