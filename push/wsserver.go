package push

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients     map[string]map[string]*Client // clinets[service][userId] = client
	addClientCh chan *Client
	rmClientCh  chan *Client
	MessageCh   chan *Message
}

func NewServer() *Server {
	return &Server{
		clients:     map[string]map[string]*Client{},
		addClientCh: make(chan *Client),
		rmClientCh:  make(chan *Client),
		MessageCh:   make(chan *Message),
	}
}

func (server *Server) addClient(client *Client) {
	fmt.Println("client add")
	_, ok := server.clients[client.Service]
	if !ok {
		server.clients[client.Service] = map[string]*Client{}
	}
	server.clients[client.Service][client.Id] = client

	// send complete message
	s := []string{"_WsProxy_REGISTERED_USER: serviceId=", client.Service, ", user=", client.Id}
	client.Send(strings.Join(s, ""))
}

func (server *Server) rmClient(client *Client) {
	fmt.Println("client remove")
	delete(server.clients[client.Service], client.Id)
}

func (server *Server) sendMessage(message *Message) {
	fmt.Println("send message")
	for _, id := range message.User {
		c, ok := server.clients[message.Service][id]
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

/**
 * 接続中のクライアントを返却するハンドラー
 * クエリー名[service]に指定したサービスで絞り込むことができます。
 * クエリー指定がない場合は、接続中の全クライアント情報を返却します。
 * ex) http://domain:port/list?service=test
 */
func (server *Server) ClientListHandler(w http.ResponseWriter, r *http.Request) {

	m := map[string][]string{}
	u := []string{}

	query := r.URL.Query()
	v, ok := query["service"]
	if ok {
		service := v[0]
		for userId, _ := range server.clients[service] {
			u = append(u, userId)
		}
		m[service] = u
	} else {
		for service, _ := range server.clients {
			u := []string{}
			for userId, _ := range server.clients[service] {
				u = append(u, userId)
			}
			m[service] = u
		}
	}

	jsonString, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonString))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)

}
