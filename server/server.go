package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Order - структура для представления данных о заказе
type Order struct {
	ID    int     `json:"id"`
	Item  string  `json:"item"`
	Price float64 `json:"price"`
	// Другие поля заказа
}

func main() {
	// Кеш для хранения данных о заказах (мапа)
	orderCache := make(map[int]Order)

	// Предварительное заполнение кеша данными (пример)
	orderCache[1] = Order{ID: 1, Item: "Laptop", Price: 1200.00}
	orderCache[2] = Order{ID: 2, Item: "Mouse", Price: 25.00}
	orderCache[3] = Order{ID: 3, Item: "Keyboard", Price: 100.00}

	// Инициализация Gin router
	router := gin.Default()

	// Эндпоинт для получения заказа по ID
	router.GET("/orders/:id", func(c *gin.Context) {
		idParam := c.Param("id") // Получаем ID из URL
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		order, exists := orderCache[id]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusOK, order)
	})

	// Запуск сервера
	port := 8070 // Порт на котором будет работать сервер
	fmt.Printf("Server is running on port: %d\n", port)
	router.Run(fmt.Sprintf(":%d", port))
}
