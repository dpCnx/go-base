package client

import (
	"log"

	"github.com/streadway/amqp"
)

func producer() {
	/*
		参数：
		1. exchange:交换机名称
		2. type:交换机类型
		FANOUT("fanout"),：扇形（广播），发送消息到每一个与之绑定队列。
		DIRECT("direct"),：定向 -->绑定key:info 发送key:info
		TOPIC("topic"),通配符的方式 -->绑定key: order.*  #.error 发送key:order.all all.error
		HEADERS("headers");参数匹配

		3. durable:是否持久化
		4. autoDelete:自动删除
		5. internal：内部使用。 一般false 	true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		6. noWait:是否阻塞处理
		7. arguments：参数
	*/
	// 尝试创建交换机
	if err := channel.ExchangeDeclare("demo", "fanout", true, false, false, false, nil); err != nil {
		log.Printf("channel ExchangeDeclare err:%v", err)
		return
	}

	// 创建队列
	q, err := channel.QueueDeclare("demo", true, false, true, false, nil)
	if err != nil {
		log.Printf("channel QueueDeclare err:%v", err)
		return
	}
	/*
		参数：
		1. queue：队列名称
		2. exchange：交换机名称
		3. routingKey：路由键，绑定规则
		如果交换机的类型为fanout ，routingKey设置为""
	*/
	// 绑定队列到 exchange
	if err = channel.QueueBind(q.Name, "", "demo", false, nil); err != nil {
		log.Printf("channel QueueBind err:%v", err)
		return
	}

	for i := 0; i <= 10; i++ {
		/*
			交换机的名字
			队列名字
			如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
			如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		*/
		if err = channel.Publish("demo", "", false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("demo"),
			}); err != nil {
			log.Println("channel Publish:", err)
			return
		}
	}
}

/*
	要注意key,规则
	其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
	匹配 dp.* 表示匹配 dp.hello, dp.hello.dd.#才能匹配到
*/
