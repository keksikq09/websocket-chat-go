package server

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
	"websocket-chat-go/src/user"
)

type Server struct {
	users map[*websocket.Conn]*user.User
	mu    sync.Mutex
}

func NewServer() *Server {
	return &Server{
		users: make(map[*websocket.Conn]*user.User),
	}
}

func (s *Server) HandleWebSocket() websocket.Handler {
	return websocket.Handler(s.handleConnection)
}

func (s *Server) handleConnection(ws *websocket.Conn) {
	user, err := s.registerUser(ws)
	if err != nil {
		log.Println("Ошибка регистрации пользователя:", err)
		return
	}

	s.broadcast(fmt.Sprintf("Пользователь %s присоединился к чату", user.Username), ws)
	s.handleMessages(user)
}

func (s *Server) registerUser(ws *websocket.Conn) (*user.User, error) {
	if err := websocket.Message.Send(ws, "Введите ваш ник:"); err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса ника: %w", err)
	}

	var nickname string
	if err := websocket.Message.Receive(ws, &nickname); err != nil {
		return nil, fmt.Errorf("ошибка при получении ника: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range s.users {
		if u.Username == nickname {
			websocket.Message.Send(ws, "Этот ник уже занят. Попробуйте другой.")
			return nil, fmt.Errorf("ник %s уже занят", nickname)
		}
	}

	newUser := &user.User{Username: nickname, Socket: ws}
	s.users[ws] = newUser

	if err := websocket.Message.Send(ws, "Добро пожаловать в чат!"); err != nil {
		return nil, fmt.Errorf("ошибка при отправке подтверждения: %w", err)
	}

	return newUser, nil
}

func (s *Server) handleMessages(user *user.User) {
	for {
		var message string
		if err := websocket.Message.Receive(user.Socket, &message); err != nil {
			log.Printf("Ошибка при получении сообщения от %s: %v", user.Username, err)
			s.removeUser(user)
			s.broadcast(fmt.Sprintf("Пользователь %s покинул чат", user.Username), user.Socket)
			return
		}

		if message[0] == '/' {
			s.handleCommand(message, user)
		} else {
			s.broadcast(fmt.Sprintf("%s: %s", user.Username, message), user.Socket)
		}
	}
}

func (s *Server) handleCommand(message string, sender *user.User) {
	parts := strings.Fields(message)

	if parts[0] == "/pm" {
		var user user.User

		for _, u := range s.users {
			if u.Username == parts[1] {
				user = *u
			}
		}

		messageContent := strings.Join(parts[2:], " ")

		s.sendToUser(messageContent, &user)

	}
}

func (s *Server) sendToUser(message string, reciver *user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ws, u := range s.users {
		if ws == reciver.Socket {
			if err := websocket.Message.Send(ws, message); err != nil {
				log.Printf("Ошибка отправки сообщения для %s: %v", u.Username, err)
			}
		}
	}
}

func (s *Server) broadcast(message string, sender *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ws, u := range s.users {
		if ws != sender {
			if err := websocket.Message.Send(ws, message); err != nil {
				log.Printf("Ошибка отправки сообщения для %s: %v", u.Username, err)
			}
		}
	}
}

func (s *Server) removeUser(user *user.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, user.Socket)
}
