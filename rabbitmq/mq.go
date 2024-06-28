package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Mq struct {
	conn               *amqp.Connection
	ch                 *amqp.Channel
	closed             bool
	closeCtx           context.Context
	closeCtxCancelFunc context.CancelFunc
}

func NewMq(url string) (*Mq, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Mq{
		conn:               conn,
		ch:                 channel,
		closeCtx:           ctx,
		closeCtxCancelFunc: cancelFunc,
	}, nil
}

func (mq *Mq) Close() error {
	if mq.closed {
		return nil
	}
	if err := mq.ch.Close(); err != nil {
		return err
	}
	if err := mq.conn.Close(); err != nil {
		return err
	}
	mq.closed = true
	mq.closeCtxCancelFunc()
	return nil
}

func (mq *Mq) CreateExchange(name string, kind string) error {
	/*
		durable:是否持久化
		autoDelete:自动删除
		internal：内部使用。 一般false 	true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		noWait:是否阻塞处理
		arguments：参数
	*/
	return mq.ch.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

func (mq *Mq) CreateQueue(name string) error {
	/*
		exclusive:独占队列（当前声明队列的连接关闭后即被删除）
	*/
	_, err := mq.ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (mq *Mq) BindQueue(queueName string, key string, exchangeName string) error {
	return mq.ch.QueueBind(queueName, key, exchangeName, false, nil)
}

func (mq *Mq) ConfirmCallBack() (chan amqp.Confirmation, error) {
	if err := mq.ch.Confirm(false); err != nil {
		return nil, err
	}
	notifyChan := make(chan amqp.Confirmation)
	mq.ch.NotifyPublish(notifyChan)
	return notifyChan, nil
}

func (mq *Mq) ReturnCallBack() chan amqp.Return {
	returnChan := make(chan amqp.Return)
	mq.ch.NotifyReturn(returnChan)
	return returnChan
}

func (mq *Mq) Producer(exchangeName string, key string) *Producer {
	return NewProducer(mq.ch, exchangeName, key)
}

func (mq *Mq) Consumer(queueName string, key string) (*Consumer, error) {
	consumer, err := NewConsumer(mq.ch, queueName, key)
	if err != nil {
		return nil, err
	}
	if err = consumer.getMsgFromCh(mq.closeCtx); err != nil {
		return nil, err
	}
	return consumer, nil
}
