package main

import (
	"net/http"
	"webpush/push"
)

func main() {
	server := push.NewServer()
	go server.Start()

	http.Handle("/ws", server.WebsocketHandler())
	http.ListenAndServe(":8248", nil)
}
