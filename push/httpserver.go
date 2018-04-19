package push

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpServer struct {
	messageCh chan *Message
}

type Message struct {
	Service string
	User    []string
	Message string
}

func NewHttpServer(message chan *Message) *HttpServer {
	return &HttpServer{
		messageCh: message,
	}
}

func (httpserver *HttpServer) Handler(w http.ResponseWriter, r *http.Request) {
	decorder := json.NewDecoder(r.Body)
	var m Message
	err := decorder.Decode(&m)
	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	fmt.Println(m.Message)
	fmt.Println(m.Service)
	fmt.Println(m.User)

	httpserver.messageCh <- &m
}

func (httpserver *HttpServer) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
