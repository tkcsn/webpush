package push

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients     map[float64]*Client
	addClientCh chan *Client
	rmClientCh  chan *Client
	MessageCh   chan *Message
}

func NewServer() *Server {
	return &Server{
		clients:     map[float64]*Client{},
		addClientCh: make(chan *Client),
		rmClientCh:  make(chan *Client),
		MessageCh:   make(chan *Message),
	}
}

func (server *Server) addClient(client *Client) {
	fmt.Println("client add")
	server.clients[client.Id] = client
}

func (server *Server) rmClient(client *Client) {
	fmt.Println("client remove")
	delete(server.clients, client.Id)
}

func (server *Server) sendMessage(message *Message) {
	fmt.Println("send message")
	for _, id := range message.User {
		c, ok := server.clients[id]
		if ok {
			go func() { c.Send(message.Message) }()
		}
	}
}

func (server *Server) Start() {
	for {
		select {
		case client := <-server.addClientCh:
			server.addClient(client)
		case client := <-server.rmClientCh:
			server.rmClient(client)
		case message := <-server.MessageCh:
			server.sendMessage(message)
		}
	}
}

func (server *Server) WebsocketHandler() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		client := NewClient(ws, server.addClientCh, server.rmClientCh)
		client.Start()
	})
}
