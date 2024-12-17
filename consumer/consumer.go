package consumer

import (
	"context"
	"database/sql"
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
		return fmt.Errorf("ошибка создания Consumer Group: %w", err)
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

	log.Println("Kafka Consumer запущен и ожидает сообщения...")
	<-ctx.Done()
	log.Println("Kafka Consumer остановлен.")
	return nil
}

type ConsumerGroupHandler struct {
	Cache       map[string]producer.Order
	UpdateCache func(order producer.Order)
	DB          *sql.DB
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
		h.processOrder(session, order)
		log.Printf("Получено сообщение: %v ", order)
		session.MarkMessage(message, "")

	}
	return nil
}

func (h *ConsumerGroupHandler) processOrder(session sarama.ConsumerGroupSession, order producer.Order) {
	ctx := context.Background()
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка начала транзакции: %v", err)
		return
	}

	// Запись в таблицу orders
	_, err = tx.ExecContext(ctx, `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи в orders: %v", err)
		return
	}

	// Запись в таблицу deliveries
	_, err = tx.ExecContext(ctx, `
        INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи в deliveries: %v", err)
		return
	}

	// Запись в таблицу payments
	_, err = tx.ExecContext(ctx, `
        INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		log.Printf("Ошибка при записи в payments: %v", err)
		return
	}

	// Запись в таблицу items
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			log.Printf("Ошибка при записи в items: %v", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Ошибка коммита транзакции: %v", err)
		return
	}
}
