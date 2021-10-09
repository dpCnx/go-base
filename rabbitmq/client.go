package client

import (
	"log"

	"github.com/streadway/amqp"
)

var channel *amqp.Channel

func init() {

	conn, err := amqp.Dial("amqp://dp:dp@192.168.172.128:5673/d")
	if err != nil {
		log.Println("amqp conn err:", err)
		return
	}
	channel, err = conn.Channel()
	// defer c.Close()
	if err != nil {
		log.Println("conn channel err:", err)
		return
	}

	log.Println("初始化channel successful")

	/*
		Confirm
	*/

	go func() {

		if err = channel.Confirm(false); err != nil {
			log.Println("c confirm err:", err)
			return
		}

		confirChan := channel.NotifyPublish(make(chan amqp.Confirmation))
		for cc := range confirChan {
			if cc.Ack {
				log.Println("confirm:消息发送成功")
			} else {
				// 这里表示消息发送到mq失败,可以处理失败流程
				log.Println("confirm:消息发送失败")
			}
		}

	}()

	/*
		return
	*/

	go func() {

		cr := channel.NotifyReturn(make(chan amqp.Return))
		for r := range cr {
			log.Println(string(r.Body))
		}
	}()
}
