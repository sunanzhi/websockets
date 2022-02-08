package websockets

const (
	TypeBroadcast = "type_broadcast"
	TypeOneOnOne = "type_one_on_one"
	TypeOneToMany = "type_one_to_many"
	TypeManyToMany = "type_many_to_many"
)

type SubscriberMessage struct {
	Type string
	Body SubscriberBody
}

type SubscriberBody struct {
	Protocol  int64
	Data      interface{}
	Condition interface{}
}

// Subscriber 订阅器接口
type Subscriber interface {
	Conn()
	Destroy()
	Producer(c *Client, message SubscriberMessage)
	Consumer(server *Server)
}
