package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"l0/producer"
)

type KafkaConsumer struct {
	brokers      []string
	topic        string
	group        string
	GroupHandler *ConsumerGroupHandler
}

func NewKafkaConsumer(brokers []string, topic, group string) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topic:   topic,
		group:   group,
	}
}

func (kc *KafkaConsumer) Start() error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(kc.brokers, kc.group, config)
	if err != nil {
		return fmt.Errorf("ошибка создания consumer group: %w", err)
	}
	handler := kc.GroupHandler

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{kc.topic}, handler); err != nil {
				log.Printf("Ошибка во время чтения сообщений: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("Kafka consumer запущен и ожидает сообщения")
	<-ctx.Done()
	log.Println("Kafka consumer остановлен")
	return nil
}

type ConsumerGroupHandler struct {
	Cache       map[string]producer.Order
	UpdateCache func(order producer.Order)
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var order producer.Order
		if err := json.Unmarshal(message.Value, &order); err != nil {
			log.Printf("Ошибка десериализации сообщения: %v", err)
			continue
		}

		h.UpdateCache(order)
		log.Printf("Получено сообщение: %v", order)
		session.MarkMessage(message, "")
	}
	return nil
}
