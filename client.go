package websockets

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id            string
	Send          chan []byte
	IsServer      bool
	CloseCallBack func()
	groups        map[string]*Manager
	Server        *Server
	socket        *websocket.Conn
}

func Init(id string, server *Server, conn *websocket.Conn) *Client {
	client := &Client{
		Id:       id,
		Send:     make(chan []byte),
		IsServer: true,
		Server:   server,
		socket:   conn,
	}
	go client.Write()
	go client.Read()

	return client
}

// Write: Send a message to the client
func (c *Client) Write() {
	defer func() {
		c.IsServer = false
		c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}

// Read: Process messages sent by the client
func (c *Client) Read() {
	defer func() {
		if err := recover(); err != nil {
			// 记录日志 并且做对应的回调 @todo
			fmt.Println("im recover", err)
			go c.Read()
		} else {
			// 用户正常退出
			c.IsServer = false
			c.socket.Close()
		}
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		err = c.Server.Protocol.Receive(c, message)
		if err != nil {
			return
		}
	}
}

func (c *Client)ToOne(toId string, msg interface{}) bool {
	//c.Server.Protocol.
	return true
}

func (c *Client)ToGroup(groupId string, msg interface{}) bool {
	return true
}