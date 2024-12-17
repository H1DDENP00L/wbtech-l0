package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"l0/consts"
	"l0/consumer"
	"l0/db"
	"l0/producer"
	"l0/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func main() {

	pg, err := sql.Open(db.DbDriver, db.DbSource)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	err = db.DatabaseUp(pg)
	if err != nil {
		log.Fatal("Ошибка запуска миграций", err)
	}

	defer pg.Close()

	err = pg.Ping()
	if err != nil {
		log.Fatalf("Не удалось пингануть БД: %v", err)
	}
	log.Println("Подклюение к БД успешно!")

	db.LoadCacheFromDB(pg)
	//fmt.Println(consts.OrderCache) использовал для дебага)))

	brokers := []string{"localhost:9093"}
	topic := "orders"
	group := "order-group"

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	// Инициализация Kafka Producer + Генерация сообщений + Отправка
	kafkaProducer, err := producer.NewKafkaProducer(brokers, topic, config)
	if err != nil {
		log.Fatalf("Ошибка создания Kafka Producer: %v", err)
	}
	defer kafkaProducer.Close()

	ctx, cancel := context.WithCancel(context.Background())

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

	// Инициализация Kafka Consumer
	consumerGroupHandler := &consumer.ConsumerGroupHandler{
		Cache:       consts.OrderCache,
		UpdateCache: db.UpdateOrderCache,
		DB:          pg,
	}
	kafkaConsumer := consumer.NewKafkaConsumer(brokers, topic, group)
	kafkaConsumer.GroupHandler = consumerGroupHandler

	go func() {
		if err := kafkaConsumer.Start(); err != nil {
			log.Fatalf("Ошибка запуска Kafka Consumer: %v", err)
		}
	}()

	server.RunServer()

	select {
	case <-ctx.Done():
		log.Println("Kafka Producer остановлен...")
		return
	}

}
