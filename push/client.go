package push

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Client struct {
	Id        int
	ws        *websocket.Conn
	messageCh chan string
}

func NewClient(ws *websocket.Conn, message chan string) *Client {
	return &Client{
		ws:        ws,
		messageCh: message,
	}
}

func (client *Client) Start() {
	for {
		var message string
		err := websocket.Message.Receive(client.ws, &message)
		if err != nil {
			// remove

		} else {
			// json action:regist, service:collie/shuffle, user:***
			// user regist
			client.messageCh <- message
		}
	}
}

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
