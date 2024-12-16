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
	Items []struct {
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
	} `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OofShard          string `json:"oof_shard"`
}

func GenerateOrder() Order {

	return Order{
		OrderUID:    uuid.New().String(),
		TrackNumber: "TRACK-67890",
		Entry:       "web",
		Locale:      "ru",
		CustomerID:  "customer-98765",
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
			Name:    "Иван Иванов",
			Phone:   "+79998887766",
			Zip:     "123456",
			City:    "Москва",
			Address: "ул. Ленина, д.1",
			Region:  "Московская область",
			Email:   "ivanov@example.com",
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
			Transaction:  "trans-12345",
			RequestID:    "req-67890",
			Currency:     "RUB",
			Provider:     "Visa",
			Amount:       15000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "Сбербанк",
			DeliveryCost: 300,
			GoodsTotal:   14700,
			CustomFee:    0,
		},
		Items: []struct {
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
		}{
			{
				ChrtID:      123,
				TrackNumber: "TRACK-67890",
				Price:       5000,
				Rid:         "RID-1",
				Name:        "Товар 1",
				Sale:        0,
				Size:        "M",
				TotalPrice:  5000,
				NmID:        111,
				Brand:       "Бренд 1",
				Status:      1,
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
