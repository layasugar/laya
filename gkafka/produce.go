package gkafka

import (
	"github.com/Shopify/sarama"
	"log"
)

func newProducer(kc *KafkaConfig) (sarama.SyncProducer, error) {
	config, err := getSaramaConfig(kc)
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducer(kc.Brokers, config)
	return producer, err
}

func prepareMessage(topic, message string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

func (kc *Engine) SendMsg(topic, message string) (partition int32, offset int64, err error) {
	producer, err := newProducer(kc.config)
	if err != nil {
		log.Printf("Could not create producer: %s", err)
	}
	msg := prepareMessage(topic, message)
	partition, offset, err = producer.SendMessage(msg)
	return
}
