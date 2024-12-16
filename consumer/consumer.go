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
	brokers []string
	topic   string
	group   string
	cache   map[string]producer.Order
}

func NewKafkaConsumer(brokers []string, topic, group string) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		topic:   topic,
		group:   group,
		cache:   make(map[string]producer.Order),
	}
}

// GetCache возвращает копию текущего кеша заказов
func (kc *KafkaConsumer) GetCache() map[string]producer.Order {
	cacheCopy := make(map[string]producer.Order)
	for key, value := range kc.cache {
		cacheCopy[key] = value
	}
	return cacheCopy
}

func (kc *KafkaConsumer) Start() error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(kc.brokers, kc.group, config)
	if err != nil {
		return fmt.Errorf("ошибка создания consumer group: %w", err)
	}
	handler := &consumerGroupHandler{cache: kc.cache}

	// Создаем контекст, который можно завершить вручную
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов для корректного завершения работы
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
			// Завершаем, если контекст завершен
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("Kafka consumer запущен и ожидает сообщения")
	<-ctx.Done() // Ожидание завершения контекста
	log.Println("Kafka consumer остановлен")
	return nil
}

type consumerGroupHandler struct {
	cache map[string]producer.Order
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var order producer.Order
		if err := json.Unmarshal(message.Value, &order); err != nil {
			log.Printf("Ошибка десериализации сообщения: %v", err)
			continue
		}
		h.cache[order.OrderUID] = order // Кэшируем по OrderUID
		log.Printf("Получено сообщение: %v", order)
		session.MarkMessage(message, "")
	}
	return nil
}
