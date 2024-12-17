package db

import (
	"context"
	"database/sql"
	"l0/consts"
	"l0/producer"
	"log"
	"time"
)

func LoadCacheFromDB(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, `SELECT 
        o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
        d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
        p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
    FROM orders o
    LEFT JOIN deliveries d ON o.order_uid = d.order_uid
    LEFT JOIN payments p ON o.order_uid = p.order_uid
    `)
	if err != nil {
		log.Printf("Ошибка запроса к таблицам orders, deliveries, payments: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order producer.Order

		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
		)

		if err != nil {
			log.Printf("Ошибка сканирования таблиц order, delivery, payment rows: %v", err)
			continue
		}

		items, err := loadItemsFromDB(ctx, db, order.OrderUID)
		if err != nil {
			log.Printf("Не удалой загрузить items для order'а %s: %v", order.OrderUID, err)
			continue
		}
		order.Items = items

		consts.CacheMutex.Lock()
		consts.OrderCache[order.OrderUID] = order
		consts.CacheMutex.Unlock()
	}

	if err = rows.Err(); err != nil {
		log.Printf("Не удалось считать строку из запроса: %v", err)
	}

	log.Println("Кэш успешно восстановлен из БД!!!.")
}
func loadItemsFromDB(ctx context.Context, db *sql.DB, orderUID string) ([]producer.Items, error) {
	items := []producer.Items{}

	itemRows, err := db.QueryContext(ctx, `SELECT 
	chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items WHERE order_uid = $1`, orderUID)

	if err != nil {
		return nil, err
	}

	defer itemRows.Close()

	for itemRows.Next() {
		var item producer.Items
		err = itemRows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = itemRows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func GetOrderCache() map[string]producer.Order {
	consts.CacheMutex.RLock()
	defer consts.CacheMutex.RUnlock()

	cacheCopy := make(map[string]producer.Order)
	for key, value := range consts.OrderCache {
		cacheCopy[key] = value
	}
	return cacheCopy
}

func UpdateOrderCache(order producer.Order) {
	consts.CacheMutex.Lock()
	defer consts.CacheMutex.Unlock()
	consts.OrderCache[order.OrderUID] = order

}
