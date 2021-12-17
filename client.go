package websocket

import "github.com/gorilla/websocket"

// 用户客户端
type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}