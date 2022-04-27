package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/shopping-service/config"
	"github.com/tengkuroman/microshop/shopping-service/controllers"
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

	// Routes (health check)
	r.GET("/check", controllers.HealthCheck)

	// Buyer route
	r.POST("/cart", controllers.AddProductToCart)
	r.GET("/cart", controllers.GetCartItems)
	r.PATCH("/cart", controllers.UpdateCartItem)
	r.DELETE("/cart", controllers.DropCart)
	r.GET("/cart/checkout", controllers.Checkout)

	// Run router
	r.Run()
}
