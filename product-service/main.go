package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/product-service/config"
	"github.com/tengkuroman/microshop/product-service/controllers"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

func routeNonAuth(key string, value interface{}) http.Handler {
	r := gin.Default()

	// Set allow CORS
	r.Use(cors.Default())

	// Set context
	r.Use(func(c *gin.Context) {
		c.Set(key, value)
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

	return r
}

func routeAuth(key string, value interface{}) http.Handler {
	r := gin.Default()

	// Set allow CORS
	r.Use(cors.Default())

	// Set context
	r.Use(func(c *gin.Context) {
		c.Set(key, value)
	})

	// Seller
	r.POST("/product", controllers.PostProduct)
	r.PATCH("/product/:product_id", controllers.UpdateProduct)
	r.DELETE("/product/:product_id", controllers.DeleteProduct)

	// Admin
	r.POST("/category", controllers.PostCategory)
	r.PATCH("/category/:category_id", controllers.UpdateCategory)
	r.DELETE("/category/:category_id", controllers.DeleteCategory)

	return r
}

func main() {
	// Connect database
	db := config.ConnectDatabase()
	databaseSQL, _ := db.DB()
	defer databaseSQL.Close()

	serverNonAuth := &http.Server{
		Addr:    ":8080",
		Handler: routeNonAuth("db", db),
	}

	serverAuth := &http.Server{
		Addr:    ":8081",
		Handler: routeAuth("db", db),
	}

	g.Go(func() error {
		err := serverNonAuth.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	g.Go(func() error {
		err := serverAuth.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
