package main

import (
	"net/http"

	"github.com/tkcsn/webpush/push"
)

func main() {

	wsserver := push.NewServer()
	go wsserver.Start()

	htserver := push.NewHttpServer(wsserver.MessageCh)

	ws := http.NewServeMux()
	ws.Handle("/ws", wsserver.WebsocketHandler())

	ht := http.NewServeMux()
	ht.HandleFunc("/ht", htserver.Handler)

	go func() {
		http.ListenAndServe(":8248", ws)
	}()

	http.ListenAndServe(":8249", ht)

}
