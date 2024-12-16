package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	// Конфигурация Sarama
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Адреса брокеров Kafka
	brokers := []string{"localhost:9093"} // Замените на ваши адреса брокеров

	// Создание consumer'а
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Не удалось создать consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Ошибка при закрытии consumer'а: %v", err)
		}
	}()

	// Топик для чтения сообщений
	topic := "orders"

	// Получение партиций топика
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Не удалось получить партиции для топика %s: %v", topic, err)
	}

	// Обработка сообщений из всех партиций
	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Не удалось создать consumer для партиции %d: %v", partition, err)
		}
		defer func(pc sarama.PartitionConsumer) {
			if err := pc.Close(); err != nil {
				log.Printf("Ошибка при закрытии partition consumer'а: %v", err)
			}
		}(partitionConsumer)

		// Обработка сообщений в отдельной горутине
		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				fmt.Printf("Получено сообщение из топика %s: партиция %d, офсет %d, ключ %s, значение %s\n",
					message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
			}
		}(partitionConsumer)
	}

	// Блокируем основной поток, чтобы consumer продолжал работать
	select {}
}
