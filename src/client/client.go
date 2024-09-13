package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func NewClient(url string) (*Client, error) {
	conn, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Run() error {
	defer c.conn.Close()

	fmt.Println("Welcome to the client")

	if err := c.register(); err != nil {
		return err
	}

	go c.receiveMessages()

	return c.sendMessages()
}

func (c *Client) register() error {
	var nickRequest string
	if err := websocket.Message.Receive(c.conn, &nickRequest); err != nil {
		return fmt.Errorf("failed to receive nick request: %w", err)
	}
	fmt.Println(nickRequest)

	reader := bufio.NewReader(os.Stdin)
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)
	if err := websocket.Message.Send(c.conn, nickname); err != nil {
		return fmt.Errorf("error sending nickname: %w", err)
	}

	var confirmation string
	if err := websocket.Message.Receive(c.conn, &confirmation); err != nil {
		return fmt.Errorf("failed to receive confirmation: %w", err)
	}
	fmt.Println(confirmation)

	return nil
}

func (c *Client) receiveMessages() {
	for {
		var msg string
		err := websocket.Message.Receive(c.conn, &msg)
		if err != nil {
			fmt.Println("Failed to receive message:", err)
			return
		}
		fmt.Println(msg)
	}
}

func (c *Client) sendMessages() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		err := websocket.Message.Send(c.conn, msg)
		if err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}
	}
	return scanner.Err()
}
