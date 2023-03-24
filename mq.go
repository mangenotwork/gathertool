package gathertool

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

// MQer 消息队列接口
type MQer interface {
	Producer(topic string, data []byte)
	Consumer(topic, channel string) []byte
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

func NewRabbit() MQer {
	return new(MQRabbitService)
}

func NewKafka() MQer {
	return new(MQKafkaService)
}

// ==============================================================================================================

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
func (m *MQNsqService) Consumer(topic, channel string) []byte {
	ch := make(chan []byte)
	//msgChan := make(chan *nsq.Message, 1024)

	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second * 3
	if len(channel) < 1 {
		channel = "mange"
	}
	consumer, err := nsq.NewConsumer(topic, channel, config)
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

	//for {
	//	select {
	//	case message := <-msgChan:
	//		ch <- message.Body
	//	}
	//}
	return <-ch
}

// ==============================================================================================================

// MQRabbitService Rabbit消息队列
type MQRabbitService struct {
}

func (m *MQRabbitService) Producer(topic string, data []byte) {}

func (m *MQRabbitService) Consumer(topic, channel string) []byte {
	return []byte{}
}

// ==============================================================================================================

// MQKafkaService Kafka消息队列
type MQKafkaService struct {
}

func (m *MQKafkaService) Producer(topic string, data []byte) {}

func (m *MQKafkaService) Consumer(topic, channel string) []byte {
	return []byte{}
}
