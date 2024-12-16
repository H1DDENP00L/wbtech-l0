package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"l0/consumer"
	"l0/db"
	"l0/producer"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

// Глобальная переменная для кеша (нужна блокировка для многопоточного доступа)
var orderCache = make(map[string]producer.Order)
var cacheMutex sync.RWMutex

// Функция для получения копии кэша
func getOrderCache() map[string]producer.Order {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	cacheCopy := make(map[string]producer.Order)
	for key, value := range orderCache {
		cacheCopy[key] = value
	}
	return cacheCopy
}

// Функция для обновления кэша
func updateOrderCache(order producer.Order) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	orderCache[order.OrderUID] = order
}

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

	// Инициализация Kafka Consumer
	consumerGroupHandler := &consumer.ConsumerGroupHandler{
		Cache:       orderCache,
		UpdateCache: updateOrderCache,
	}
	kafkaConsumer := consumer.NewKafkaConsumer(brokers, topic, group)
	kafkaConsumer.GroupHandler = consumerGroupHandler

	go func() {
		if err := kafkaConsumer.Start(); err != nil {
			log.Fatalf("Ошибка запуска Kafka consumer: %v", err)
		}
	}()

	// Инициализация Gin router
	router := gin.Default()

	// Загрузка HTML-шаблона
	tmpl := template.Must(template.ParseFiles("templates/order.html"))

	// Настройка HTML-рендерера
	router.SetHTMLTemplate(tmpl)

	// Эндпоинт для получения заказа по OrderUID
	router.GET("/orders/:order_uid", func(c *gin.Context) {
		orderUID := c.Param("order_uid")
		cache := getOrderCache()

		order, exists := cache[orderUID]
		if !exists {
			c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Order not found"})
			return
		}
		c.HTML(http.StatusOK, "order.html", gin.H{"order": order})
	})

	// Запуск сервера
	port := 8070 // Порт на котором будет работать сервер
	fmt.Printf("Server is running on port: %d\n", port)
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatalf("Ошибка запуска http сервера: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Kafka-Producer остановлен...")
		return
	}

}
