package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type User struct {
	Username string
	Socket   *websocket.Conn
}

var users = make(map[*websocket.Conn]*User)

func main() {
	http.Handle("/echo", websocket.Handler(handleConnection))
	fmt.Println("Server start running")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnection(ws *websocket.Conn) {
	// receive nick
	if err := websocket.Message.Send(ws, "Введите ваш ник:"); err != nil {
		log.Println("Ошибка при отправке запроса ника:", err)
		return
	}

	// get nick
	var nickname string
	if err := websocket.Message.Receive(ws, &nickname); err != nil {
		log.Println("Ошибка при получении ника:", err)
		return
	}

	// check if nickname is taken
	for _, user := range users {
		if user.Username == nickname {
			websocket.Message.Send(ws, "Этот ник уже занят. Попробуйте другой.")
			return
		}
	}

	// create new user
	user := &User{Username: nickname, Socket: ws}
	users[ws] = user

	// send confirmation
	if err := websocket.Message.Send(ws, "Добро пожаловать в чат!"); err != nil {
		log.Println("Ошибка при отправке подтверждения:", err)
		return
	}

	broadcast(fmt.Sprintf("Пользователь %s присоединился к чату", nickname), ws)
	handleMessages(user)
}

func handleMessages(user *User) {
	for {
		var message string
		if err := websocket.Message.Receive(user.Socket, &message); err != nil {
			log.Printf("Ошибка при получении сообщения от %s: %v", user.Username, err)
			delete(users, user.Socket)
			broadcast(fmt.Sprintf("Пользователь %s покинул чат", user.Username), user.Socket)
			return
		}
		broadcast(fmt.Sprintf("%s: %s", user.Username, message), user.Socket)
	}
}

func broadcast(message string, sender *websocket.Conn) {
	for ws, user := range users {
		if ws != sender { // Не отправляем сообщение отправителю
			if err := websocket.Message.Send(ws, message); err != nil {
				log.Printf("Ошибка отправки сообщения для %s: %v", user.Username, err)
			}
		}
	}
}
