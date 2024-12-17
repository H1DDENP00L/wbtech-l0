package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"l0/db"
	"log"
	"net/http"
)

func RunServer() {
	router := gin.Default()

	tmpl := template.Must(template.ParseFiles("templates/order.html"))

	router.SetHTMLTemplate(tmpl)

	// Получаем данные о заказе по OrderUID
	router.GET("/orders/:order_uid", func(c *gin.Context) {
		orderUID := c.Param("order_uid")
		cache := db.GetOrderCache()

		order, exists := cache[orderUID]
		if !exists {
			c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Order not found"})
			return
		}
		c.HTML(http.StatusOK, "order.html", gin.H{"order": order})
	})

	port := 8070
	log.Printf("Сервер запущен на порту: %d\n", port)
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatalf("Ошибка запуска http сервера: %v", err)
		}
	}()
}
