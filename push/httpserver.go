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
	service string
	User    []float64
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
	fmt.Println(m.User)

	httpserver.messageCh <- &m
}
