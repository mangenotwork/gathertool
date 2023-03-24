/*
	Description : 消息接口与相关方法, 支持 Nsq, Rabbit, Kafka
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/nsqio/go-nsq"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

// MQer 消息队列接口
type MQer interface {
	Producer(topic string, data []byte)
	Consumer(topic string) []byte
}

// NewNsq port 依次是 ProducerPort, ConsumerPort
func NewNsq(server string, port ...int) MQer {
	mq := &MQNsqService{
		NsqServerIP: server,
	}
	if len(port) > 0 {
		mq.ProducerPort = port[0]
	}
	if len(port) > 1 {
		mq.ProducerPort = port[1]
	}
	return mq
}

func NewRabbit(amqpUrl string) MQer {
	return &MQRabbitService{
		AmqpUrl: amqpUrl,
	}
}

func NewKafka(server []string) MQer {
	return &MQKafkaService{
		Server: server,
	}
}

// nsq ==============================================================================================================

// MQNsqService NSQ消息队列
type MQNsqService struct {
	NsqServerIP  string
	ProducerPort int
	ConsumerPort int
}

// Producer 生产者
func (m *MQNsqService) Producer(topic string, data []byte) {
	nsqConf := nsq.NewConfig()

	addr := m.NsqServerIP + ":4150"
	if m.ProducerPort > 0 {
		addr = fmt.Sprintf("%s:%d", m.NsqServerIP, m.ProducerPort)
	}

	client, err := nsq.NewProducer(addr, nsqConf)
	if err != nil {
		Error("[nsq]无法连接到队列")
		return
	}
	client.SetLogger(nil, 0)
	Info(fmt.Sprintf("[生产消息] topic : %s -->  %s", topic, string(data)))
	err = client.Publish(topic, data)
	if err != nil {
		Error("[生产消息] 失败 ： " + err.Error())
	}
}

// Consumer 消费者
func (m *MQNsqService) Consumer(topic string) []byte {
	ch := make(chan []byte)
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second * 3
	consumer, err := nsq.NewConsumer(topic, "mange", config)
	if err != nil {
		panic(err)
	}
	consumer.SetLogger(nil, 0)

	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		Infof("[NSQ]消费消息:%+v", m)
		ch <- m.Body
		m.Finish()
		return nil
	}))
	addr := m.NsqServerIP + ":4161"
	if m.ConsumerPort > 0 {
		addr = fmt.Sprintf("%s:%d", m.NsqServerIP, m.ConsumerPort)
	}
	if err = consumer.ConnectToNSQLookupd(addr); err != nil {
		panic(err)
	}
	return <-ch
}

// Rabbit ===========================================================================================================

// MQRabbitService Rabbit消息队列
type MQRabbitService struct {
	AmqpUrl string
}

func (m *MQRabbitService) Producer(topic string, data []byte) {
	mq, err := NewRabbitMQPubSub(topic, m.AmqpUrl)
	if err != nil {
		Error("[rabbit]无法连接到队列")
		return
	}
	//defer mq.Destroy()
	Infof(fmt.Sprintf("[生产消息] topic : %s -->  %s", topic, string(data)))
	err = mq.PublishPub(data)
	if err != nil {
		Error("[生产消息] 失败 ： " + err.Error())
	}
}

func (m *MQRabbitService) Consumer(topic string) []byte {
	mh, err := NewRabbitMQPubSub(topic, m.AmqpUrl)
	if err != nil {
		Error("[rabbit]无法连接到队列")
		return []byte{}
	}
	msg := mh.RegistryReceiveSub()
	s := <-msg
	return s.Body
}

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string //交换机
	Key       string //key Simple模式 几乎用不到
	MqUrl     string //连接信息
}

// NewRabbitMQ 创建RabbitMQ结构体实例
func NewRabbitMQ(queueName, exchange, key, amqpUrl string) (*RabbitMQ, error) {
	mq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: amqpUrl}
	var err error
	//创建rabbitMq连接
	mq.conn, err = amqp.Dial(mq.MqUrl)
	if err != nil {
		return nil, err
	}
	mq.failOnErr(err, "创建连接错误！")
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "获取channel失败")
	return mq, err
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		Error(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 简单模式step：1。创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) (*RabbitMQ, error) {
	return NewRabbitMQ(queueName, "", "", "")
}

// PublishSimple 简单模式Step:2、简单模式下生产代码
func (r *RabbitMQ) PublishSimple(message []byte) (err error) {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err = r.channel.QueueDeclare(
		r.QueueName, //队列名称
		false,       //是否持久化
		false,       //是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,       //是否具有排他性 true表示自己可见 其他用户不能访问
		false,       //是否阻塞 true表示要等待服务器的响应
		nil,         //额外数学系
	)
	r.failOnErr(err, "failed to declare a queue")
	//2.发送消息到队列中
	err = r.channel.Publish(
		r.Exchange,  //默认的Exchange交换机是default,类型是direct直接类型
		r.QueueName, //要赋值的队列名称
		false,       //如果为true，根据exchange类型和rout key规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,       //如果为true,当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息还给发送者
		amqp.Publishing{
			ContentType: "text/plain", //类型
			Body:        message,      //消息
		})
	r.failOnErr(err, "publish 消息失败")
	return
}

// RegistryConsumeSimple 简单模式注册消费者
func (r *RabbitMQ) RegistryConsumeSimple() (msg <-chan amqp.Delivery) {
	//1、申请队列，如果队列存在就跳过，不存在创建
	//优点：保证队列存在，消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName, //队列名称
		false,       //是否持久化
		false,       //是否为自动删除 当最后一个消费者断开连接之后，是否把消息从队列中删除
		false,       //是否具有排他性
		false,       //是否阻塞
		nil,         //额外参数
	)
	if err != nil {
		fmt.Println(err)
	}
	//接收消息
	msg, err = r.channel.Consume(
		r.QueueName,
		"",    //用来区分多个消费者
		true,  //是否自动应答
		false, //是否具有排他性
		false, //如果设置为true,表示不能同一个connection中发送的消息传递给这个connection中的消费者
		false, //队列是否阻塞
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// NewRabbitMQPubSub 订阅模式创建 rabbitMq实例  (目前用的fanout模式)
func NewRabbitMQPubSub(exchangeName, amqpUrl string) (*RabbitMQ, error) {
	mq, err := NewRabbitMQ("", exchangeName, "", amqpUrl)
	if mq == nil || err != nil {
		return nil, err
	}
	//获取connection
	mq.conn, err = amqp.Dial(mq.MqUrl)
	mq.failOnErr(err, "failed to connect mq!")
	if mq.conn == nil || err != nil {
		return nil, err
	}
	//获取channel
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "failed to open a channel!")
	return mq, err
}

// PublishPub 订阅模式生成
func (r *RabbitMQ) PublishPub(message []byte) (err error) {
	//尝试创建交换机，不存在创建
	err = r.channel.ExchangeDeclare(
		r.Exchange,          //交换机名称
		amqp.ExchangeFanout, //交换机类型 广播类型
		true,                //是否持久化
		false,               //是否自动删除
		false,               //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,               //是否阻塞 true表示要等待服务器的响应  false 无等待
		nil,                 //参数
	)
	r.failOnErr(err, "failed to declare an exchange"+"nge")
	//2 发送消息
	err = r.channel.Publish(r.Exchange, "", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return
}

// RegistryReceiveSub 订阅模式消费端代码
func (r *RabbitMQ) RegistryReceiveSub() (msg <-chan amqp.Delivery) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		r.Exchange,          //交换机名称
		amqp.ExchangeFanout, //交换机类型 广播类型
		true,                //是否持久化
		false,               //是否字段删除
		false,               //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,               //是否阻塞 true表示要等待服务器的响应
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	//2. 试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	//在pub/sub模式下，这里的key要为空
	err = r.channel.QueueBind(q.Name, "", r.Exchange, false, nil)
	//消费消息
	msg, err = r.channel.Consume(q.Name, "", true, false, false, false, nil)
	return
}

// NewRabbitMQTopic 话题模式 创建RabbitMQ实例
func NewRabbitMQTopic(exchange string, routingKey string) (*RabbitMQ, error) {
	mq, _ := NewRabbitMQ("", exchange, routingKey, "")
	var err error
	mq.conn, err = amqp.Dial(mq.MqUrl)
	mq.failOnErr(err, "failed   to connect rabbitMq!")
	mq.channel, err = mq.conn.Channel()
	mq.failOnErr(err, "failed to open a channel")
	return mq, err
}

// PublishTopic 话题模式发送信息
func (r *RabbitMQ) PublishTopic(message []byte) (err error) {
	//尝试创建交换机，不存在创建
	err = r.channel.ExchangeDeclare(
		r.Exchange,         //交换机名称
		amqp.ExchangeTopic, //交换机类型 话题模式
		true,               //是否持久化
		false,              //是否字段删除
		false,              //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,              //是否阻塞 true表示要等待服务器的响应
		nil,
	)
	r.failOnErr(err, "topic failed to declare an exchange")
	//2发送信息
	err = r.channel.Publish(r.Exchange, r.Key, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return
}

// RegistryReceiveTopic 话题模式接收信息
// 要注意key
// 其中* 用于匹配一个单词，#用于匹配多个单词（可以是零个）
// 匹配 xx.* 表示匹配xx.hello,但是xx.hello.one需要用xx.#才能匹配到
func (r *RabbitMQ) RegistryReceiveTopic() (msg <-chan amqp.Delivery) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		r.Exchange,         //交换机名称
		amqp.ExchangeTopic, //交换机类型 话题模式
		true,               //是否持久化
		false,              //是否字段删除
		false,              //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,              //是否阻塞 true表示要等待服务器的响应
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	//2. 试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	//在pub/sub模式下，这里的key要为空
	err = r.channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
	//消费消息
	msg, err = r.channel.Consume(q.Name, "", true, false, false, false, nil)
	return
}

// NewRabbitMQRouting 路由模式 创建RabbitMQ实例
func NewRabbitMQRouting(exchange string, routingKey string) (*RabbitMQ, error) {
	rabbitMQ, _ := NewRabbitMQ("", exchange, routingKey, "")
	var err error
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.MqUrl)
	rabbitMQ.failOnErr(err, "failed   to connect rabbitMq!")
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open a channel")
	return rabbitMQ, err
}

// PublishRouting 路由模式发送信息
func (r *RabbitMQ) PublishRouting(message []byte) (err error) {
	//尝试创建交换机，不存在创建
	err = r.channel.ExchangeDeclare(
		r.Exchange,          //交换机名称
		amqp.ExchangeDirect, //交换机类型 广播类型
		true,                //是否持久化
		false,               //是否字段删除
		false,               //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,               //是否阻塞 true表示要等待服务器的响应
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")
	//发送信息
	err = r.channel.Publish(r.Exchange, r.Key, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return
}

// RegistryReceiveRouting 路由模式接收信息
func (r *RabbitMQ) RegistryReceiveRouting() (msg <-chan amqp.Delivery) {
	//尝试创建交换机，不存在创建
	err := r.channel.ExchangeDeclare(
		r.Exchange,          //交换机名称
		amqp.ExchangeDirect, //交换机类型 广播类型
		true,                //是否持久化
		false,               //是否字段删除
		false,               //true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,               //是否阻塞 true表示要等待服务器的响应
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange"+"nge")
	// 试探性创建队列，创建队列
	q, err := r.channel.QueueDeclare("", false, false, true, false, nil)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到exchange中
	//在pub/sub模式下，这里的key要为空
	err = r.channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
	//消费消息
	msg, err = r.channel.Consume(q.Name, "", true, false, false, false, nil)
	return
}

// Kafka ===========================================================================================================

// MQKafkaService Kafka消息队列
type MQKafkaService struct {
	Server []string
}

func (m *MQKafkaService) Producer(topic string, data []byte) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //写到随机分区中，我们默认设置32个分区
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.ByteEncoder(data)
	// 连接kafka
	client, err := sarama.NewSyncProducer(m.Server, config)
	if err != nil {
		Error("Producer closed, err:", err)
		return
	}
	defer func() {
		_ = client.Close()
	}()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		Error("send msg failed, err:", err)
		return
	}
	Infof("pid:%v offset:%v\n", pid, offset)
}

func (m *MQKafkaService) Consumer(topic string) []byte {
	ch := make(chan []byte)
	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer(m.Server, nil)
	if err != nil {
		Errorf("Failed to start consumer: %s", err)
		return []byte{}
	}
	partitionList, err := consumer.Partitions("task-status-data") // 通过topic获取到所有的分区
	if err != nil {
		Error("Failed to get the list of partition: ", err)
		return []byte{}
	}
	Info(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest) // 针对每个分区创建一个分区消费者
		if err != nil {
			Errorf("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { // 为每个分区开一个go协程取值
			for msg := range pc.Messages() { // 阻塞直到有值发送过来，然后再继续等待
				Infof("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				ch <- msg.Value
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()

	_ = consumer.Close()

	return <-ch
}
