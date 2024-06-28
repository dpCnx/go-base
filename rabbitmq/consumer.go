package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

const msgChannelSize = 1000

type Consumer struct {
	ch        *amqp.Channel
	queueName string
	key       string
	msgCh     chan []byte
}

func NewConsumer(ch *amqp.Channel, queueName string, key string) (*Consumer, error) {

	/*
		能者多劳,每次只能获取一条消息,处理完成才能获取下一个消息
		prefetch count
		prefetch size in bytes; 0 is ignored
		global prefetch 计数（或大小）是应用于单个消费者还是所有消费者 一般都设置为false
	*/
	if err := ch.Qos(1, 0, false); err != nil {
		return nil, err
	}
	return &Consumer{
		ch:        ch,
		queueName: queueName,
		key:       key,
		msgCh:     make(chan []byte, msgChannelSize),
	}, nil
}

func (cm *Consumer) ConsumeChan() (<-chan []byte, error) {
	return cm.msgCh, nil
}

func (cm *Consumer) getMsgFromCh(ctx context.Context) error {
	/*
		queueName
		consumer:用来区分多个消费者
		autoAck:是否自动应答
		exclusive:是否独有
		noLocal:设置为true,表示 不能将同一个Connection中生产者发送的消息传递给这个Connection中 的消费者
		noWait:是否阻塞
	*/
	consumerMsgCh, err := cm.ch.Consume(cm.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-consumerMsgCh:
				cm.msgCh <- msg.Body
			/*
				autoAck 设置为false 手动签收
				告知RabbitMQ消费者已成功处理某条消息，可以将该消息从队列中删除
				msg.Ack(true) // 若设置为true，则确认所有小于或等于deliveryTag的消息；若设置为false，则仅确认deliveryTag对应的消息
				告知RabbitMQ消费者无法处理某条消息，并决定该消息的处理方式（重新入队或丢弃）
				msg.Reject(true) // 若设置为true，则将消息重新放入队列，等待其他消费者处理；若设置为false，则消息将被丢弃。

				// 若设置为true，则确认所有小于或等于deliveryTag的消息；若设置为false，则仅确认deliveryTag对应的消息
				// 若设置为true，则拒绝所有小于或等于deliveryTag的消息；若设置为false，则仅拒绝deliveryTag对应的消息。
				msg.Nack(true, true)
			*/
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
