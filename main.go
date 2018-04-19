package main

import (
	"log"
	"net/http"
	"webpush/push"
)

func main() {

	wsserver := push.NewServer()
	go wsserver.Start()

	htserver := push.NewHttpServer(wsserver.MessageCh)

	ws := http.NewServeMux()
	ws.Handle("/ws", wsserver.WebsocketHandler())

	ht := http.NewServeMux()
	ht.HandleFunc("/ht", htserver.Handler)
	ht.HandleFunc("/status", htserver.HealthHandler)
	ht.HandleFunc("/list", wsserver.ClientListHandler)

	go func() {
		err := http.ListenAndServe(":8248", ws)
		if err != nil {
			log.Fatal("wsserver: ", err)
		}
	}()

	http.ListenAndServe(":8249", ht)

}
