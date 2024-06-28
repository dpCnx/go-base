package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	ch *amqp.Channel

	exchangeName string
	key          string
}

func NewProducer(ch *amqp.Channel, exchangeName string, key string) *Producer {
	return &Producer{
		ch:           ch,
		exchangeName: exchangeName,
		key:          key,
	}
}

func (p *Producer) produce(msg []byte) error {
	/*
		如果为true,根据自身exchange类型和routeKey规则无法找到符合条件的队列会把消息返还给发送者
		如果为true,当exchange发送消息到队列后发现队列上没有消费者,则会把消息返还给发送者
	*/
	return p.ch.Publish(p.exchangeName, p.key, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
}
