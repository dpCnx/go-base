package client

import (
	"log"
)

func consumer() {
	/*
		queue
		用来区分多个消费者
		是否自动应答
		是否独有
		设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		是否阻塞
	*/
	msgs, err := channel.Consume("demo", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("get msg err: %v", err)
		return
	}

	go func() {
		for d := range msgs {
			log.Printf("message: %s\n", d.Body)
			// autoAck 设置为false 手动签收
			if err := d.Ack(true); err != nil {
				log.Println(err.Error())
				return
			}
			// 重回队列。如果设置为true，则消息重新回到queue，broker会重新发送该消息给消费端
			_ = d.Nack(true, true)
		}
	}()

}
