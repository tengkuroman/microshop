package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/payment-service/config"
	"github.com/tengkuroman/microshop/payment-service/controllers"
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

	// Routes (public)
	r.GET("/payment", controllers.GetPaymentProviders)

	// Routes (admin)
	r.POST("/payment", controllers.PostPaymentProvider)
	r.PATCH("/payment/:payment_provider_id", controllers.UpdatePaymentProvider)
	r.DELETE("/payment/:payment_provider_id", controllers.DeletePaymentProvider)

	// Routes (service)
	r.POST("/payment/process", controllers.ProcessPayment)

	// Run router
	r.Run()
}
