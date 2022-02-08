package main

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"net/http"
	"websockets"
	"websockets/message"
	"websockets/parser"
)

type AddReq struct {
	Content string
}

type AddRes struct {
	Content string
}

var Server = &websockets.Server{}

func main() {
	node, _ := snowflake.NewNode(1)
	protocol := parser.Json()
	protocol.Register(1, AddReq{}, addReq)
	protocol.Register(-1, AddRes{}, nil)
	rabbitMq := message.RabbitMq("", "newProduct", "", "amqp://rabbit:rabbit@192.168.0.31:5672/test")

	Server =  websockets.InitServer(protocol, rabbitMq)
	Server.Subscriber.Conn()
	go Server.Subscriber.Consumer(Server)
	websockets.Listen(":9504", "/wss", func(res http.ResponseWriter, req *http.Request) {
		//解析一个连接
		conn, _ := (&websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}).Upgrade(res, req, nil)
		// 初始化一个客户端对象
		client := websockets.Init(node.Generate().String(), Server, conn)
		Server.Manager.Join(client)
	})
}

func addReq(c *websockets.Client, msg interface{}) {
	addReq := msg.(*AddReq)
	fmt.Println(c.Id, msg, addReq.Content, "end")
	c.Server.Subscriber.Producer(c, websockets.SubscriberMessage{
		Type: websockets.TypeBroadcast,
		Body: websockets.SubscriberBody{
			Protocol: -1,
			Data: AddRes{
				Content: "我是回复消息",
			},
		},
	})
}