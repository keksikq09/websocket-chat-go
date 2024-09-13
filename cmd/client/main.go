package main

import (
	"log"

	"websocket-chat-go/src/client"
)

func main() {
	c, err := client.NewClient("ws://localhost:8080/echo")
	if err != nil {
		log.Fatal("Error creating client:", err)
	}

	err = c.Run()
	if err != nil {
		log.Fatal("Error running client:", err)
	}
}
