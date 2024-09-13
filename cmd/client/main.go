package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"strings"
)

func main() {
	ws, err := websocket.Dial("ws://localhost:8080/echo", "", "http://localhost/")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	fmt.Println("Welcome to the client")

	// register
	var nickRequest string
	if err := websocket.Message.Receive(ws, &nickRequest); err != nil {
		log.Fatal("Failed to receive nick request:", err)
	}
	fmt.Println(nickRequest) // Выводим запрос ника

	reader := bufio.NewReader(os.Stdin)
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)
	if err := websocket.Message.Send(ws, nickname); err != nil {
		log.Fatal("Ошибка при отправке ника:", err)
	}

	// Ждем подтверждения от сервера
	var confirmation string
	if err := websocket.Message.Receive(ws, &confirmation); err != nil {
		log.Fatal("Failed to receive confirmation:", err)
	}
	fmt.Println(confirmation)

	// receive message
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println("Failed receive message:", err)
				return
			}
			fmt.Println(msg)
		}
	}()

	// send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
			return
		}
	}
}
