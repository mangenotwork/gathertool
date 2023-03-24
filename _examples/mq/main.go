package main

import (
	gt "github.com/mangenotwork/gathertool"
)

func main() {

	NsqProducer()

	//NsqConsumer()

}

func NsqProducer() {
	mq := gt.NewNsq("127.0.0.1")
	topic := "test"
	data := []byte("data")
	mq.Producer(topic, data)
}

func NsqConsumer() {
	topic := "test"
	mq := gt.NewNsq("127.0.0.1")
	for {
		data := mq.Consumer(topic)
		gt.Info(string(data))
	}

}

func RabbitProducer() {
	mq := gt.NewRabbit("amqp://admin:123456@127.0.0.1:5672")
	topic := "test"
	data := []byte("data")
	mq.Producer(topic, data)
}

func RabbitConsumer() {
	topic := "test"
	mq := gt.NewRabbit("amqp://admin:123456@127.0.0.1:5672")
	data := mq.Consumer(topic)
	gt.Info(string(data))
}
