package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/order-service/config"
	"github.com/tengkuroman/microshop/order-service/controllers"
)

func main() {
	// Connect database
	db := config.ConnectDatabase()
	databaseSQL, _ := db.DB()
	defer databaseSQL.Close()

	// Router
	r := gin.Default()

	// Set allow CORS
	r.Use(cors.Default())

	// Set context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// Routes (user)
	r.GET("/orders", controllers.GetOrdersDetail)
	r.DELETE("/order/delete/:order_detail_id", controllers.DeleteOrder)
	r.PATCH("/order/payment/:order_detail_id/:payment_provider_id", controllers.SelectPaymentProvider)
	r.PATCH("/order/payment/checkout/:order_detail_id", controllers.PayOrder)

	// Routes (service)
	r.POST("/order", controllers.CreateOrder)

	// Run router
	r.Run()
}
