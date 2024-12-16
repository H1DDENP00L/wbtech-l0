package main

import (
	"context"
	"l0/consumer"
	"l0/producer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"l0/db"
)

func main() {
	db.DatabaseUp()
	brokers := []string{"localhost:9093"}
	topic := "orders"
	group := "order-group"

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	kafkaProducer, err := producer.NewKafkaProducer(brokers, topic, config)
	if err != nil {
		log.Fatalf("Ошибка создания Kafka продюсера: %v", err)
	}
	defer kafkaProducer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	// Обработка сигналов для корректного завершения работы
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}()

	go func() {
		tik := time.NewTicker(time.Second * 5)

		for range tik.C {
			order := producer.GenerateOrder()
			if err := kafkaProducer.SendMessage(order); err != nil {
				log.Fatalf("Ошибка отправки сообщения: %v", err)
			}

		}
	}()

	KafkaConsumer := consumer.NewKafkaConsumer(brokers, topic, group)
	if err := KafkaConsumer.Start(); err != nil {
		log.Fatalf("Ошибка запуска консюмера: %v", err)
	}

	select {
	case <-ctx.Done():
		log.Println("Kafka-Producer остановлен...")
		return

	}
}
