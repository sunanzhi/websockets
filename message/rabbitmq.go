package message

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"websockets"
)

type RabbitmqMessage struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//bind Key 名称
	Key string
	//连接信息
	MqUrl string
}

func RabbitMq(queueName string, exchange string, key string, mqUrl string) websockets.Subscriber {
	return &RabbitmqMessage{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: mqUrl}
}

func (r *RabbitmqMessage) Conn() {
	//获取connection
	var err error
	r.conn, err = amqp.Dial(r.MqUrl)
	r.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	r.channel, err = r.conn.Channel()
	r.failOnErr(err, "failed to open a channel")
}

// Destroy 断开channel 和 connection
func (r *RabbitmqMessage) Destroy() {
	err := r.channel.Close()
	if err != nil {
		return
	}
	err = r.conn.Close()
	if err != nil {
		return
	}
}

// Producer 订阅模式生产
func (r *RabbitmqMessage) Producer(c *websockets.Client, message websockets.SubscriberMessage) {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	body, _ := json.Marshal(message)
	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

// Consumer 订阅模式消费端代码
func (r *RabbitmqMessage) Consumer(server *websockets.Server) {
	// 试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//交换机类型
		"fanout",
		true,
		false,
		false, // true 表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")
	// 试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		// 在pub/sub模式下，这里的key要为空
		"",
		r.Exchange,
		false,
		nil)

	//消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for d := range messages {
			message := websockets.SubscriberMessage{}
			json.Unmarshal(d.Body, &message)
			switch message.Type {
			case websockets.TypeBroadcast:
				server.Manager.Broadcast(message.Body)
				break
			}
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}



//错误处理函数
func (r *RabbitmqMessage) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
	}
}
