package user

import (
	"golang.org/x/net/websocket"
)

type User struct {
	Username string
	Socket   *websocket.Conn
}
