package websockets

type ParserHandler func(c *Client, msg interface{})

// Protocol 协议接口
type Protocol interface {
	Register(protocol int64, t interface{}, handler ParserHandler)
	Receive(c *Client, msg []byte) error
	Send(c *Client, msg SubscriberBody) error
}