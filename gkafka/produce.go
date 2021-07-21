package gkafka

import (
	"github.com/Shopify/sarama"
)

func newProducer(kc *KafkaConfig) (sarama.SyncProducer, error) {
	config, err := getSaramaConfig(kc)
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducer(kc.Brokers, config)
	return producer, err
}

func prepareMessage(topic, message string, partition int32) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: partition,
		Value:     sarama.StringEncoder(message),
	}

	return msg
}

func (kc *Engine) InitProducer(config *KafkaConfig) error {
	producer, err := newProducer(config)
	if err != nil {
		return err
	}
	kc.producer = producer
	return nil
}

func (kc *Engine) SendMsg(topic, message string, part int32) (partition int32, offset int64, err error) {
	msg := prepareMessage(topic, message, part)
	partition, offset, err = kc.producer.SendMessage(msg)
	return
}
