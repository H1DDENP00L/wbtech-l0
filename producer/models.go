package producer

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int64  `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items             []Items `json:"items"`
	Locale            string  `json:"locale"`
	InternalSignature string  `json:"internal_signature"`
	CustomerID        string  `json:"customer_id"`
	DeliveryService   string  `json:"delivery_service"`
	Shardkey          string  `json:"shardkey"`
	SmID              int     `json:"sm_id"`
	DateCreated       string  `json:"date_created"`
	OofShard          string  `json:"oof_shard"`
}

type Items struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func GenerateOrder() Order {

	return Order{
		OrderUID:    uuid.New().String(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Locale:      "ru",
		CustomerID:  "customer-46064",
		DateCreated: time.Now().Format(time.RFC3339),
		Delivery: struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		}{
			Name:    "Владимир Максимович",
			Phone:   "+79998887766",
			Zip:     "606066",
			City:    "Москва",
			Address: "ул. Большая Садовая, 10",
			Region:  "Московская область",
			Email:   "wbEnjoyer@example.com",
		},
		Payment: struct {
			Transaction  string `json:"transaction"`
			RequestID    string `json:"request_id"`
			Currency     string `json:"currency"`
			Provider     string `json:"provider"`
			Amount       int    `json:"amount"`
			PaymentDt    int64  `json:"payment_dt"`
			Bank         string `json:"bank"`
			DeliveryCost int    `json:"delivery_cost"`
			GoodsTotal   int    `json:"goods_total"`
			CustomFee    int    `json:"custom_fee"`
		}{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "req-46064",
			Currency:     "RUB",
			Provider:     "MIR",
			Amount:       1817,
			PaymentDt:    time.Now().Unix(),
			Bank:         "Sber",
			DeliveryCost: 0,
			GoodsTotal:   12054,
			CustomFee:    0,
		},
		Items: []Items{
			{
				ChrtID:      123,
				TrackNumber: "WBILMTESTTRACK",
				Price:       12054,
				Rid:         "ab5212088b777ae0btest",
				Name:        "Рубашка",
				Sale:        12,
				Size:        "L",
				TotalPrice:  12054,
				NmID:        12688120,
				Brand:       "bikkembergs",
				Status:      202,
			},
		},
	}
}

//func randomStr(str string, lenStr int) string {
//	var b strings.Builder
//	for i := 0; i < lenStr; i++ {
//		ch := string(str[rand.Intn(len(str))])
//		b.WriteString(ch)
//	}
//	return b.String()
//}
