package producer

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string, config *sarama.Config) (*KafkaProducer, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Kafka продюсера: %w", err)
	}
	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) SendMessage(order interface{}) error {
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("ошибка сериализации сообщения: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.StringEncoder(orderJSON),
	}

	partition, offset, err := kp.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}

	log.Printf("Сообщение отправлено в топик %s [partition: %d, offset: %d]\n", kp.topic, partition, offset)
	return nil
}

func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}
