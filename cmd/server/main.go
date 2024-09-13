package main

import (
	"fmt"
	"log"
	"net/http"

	"websocket-chat-go/src/server"
)

func main() {
	s := server.NewServer()

	http.Handle("/echo", s.HandleWebSocket())

	fmt.Println("Сервер запущен на :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
