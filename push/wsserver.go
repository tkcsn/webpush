package push

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients     map[int]*Client
	addClientCh chan *Client
	messageCh   chan string
}

func NewServer() *Server {
	return &Server{
		clients:     map[int]*Client{},
		addClientCh: make(chan *Client),
		messageCh:   make(chan string),
	}
}

func (server *Server) addClient(client *Client) {
	fmt.Println("client add")
	server.clients[0] = client
}

func (server *Server) sendMessage(message string) {
	for _, client := range server.clients {
		c := client
		go func() { c.Send(message) }()
	}
}

func (server *Server) Start() {
	for {
		select {
		case client := <-server.addClientCh:
			server.addClient(client)
		case message := <-server.messageCh:
			server.sendMessage(message)
		}
	}
}

func (server *Server) WebsocketHandler() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		client := NewClient(ws, server.messageCh)
		server.addClientCh <- client
		client.Start()
	})
}
