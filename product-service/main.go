package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/product-service/config"
	"github.com/tengkuroman/microshop/product-service/controllers"
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

	// All user
	r.GET("/products", controllers.GetAllProducts)
	r.GET("/product/:product_id", controllers.GetProductByID)
	r.GET("/products/seller/:user_id", controllers.GetProductsBySellerID)
	r.GET("/products/category/:category_id", controllers.GetProductsByCategoryID)

	r.GET("/categories", controllers.GetAllCategories)
	r.GET("/category/:category_id", controllers.GetCategoryByID)

	// Seller
	r.POST("/product", controllers.PostProduct)
	r.PATCH("/product/:product_id", controllers.UpdateProduct)
	r.DELETE("/product/:product_id", controllers.DeleteProduct)

	// Admin
	r.POST("/category", controllers.PostCategory)
	r.PATCH("/category/:category_id", controllers.UpdateCategory)
	r.DELETE("/category/:category_id", controllers.DeleteCategory)

	// Run router
	r.Run()
}
