package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/shopping-service/config"
	"github.com/tengkuroman/microshop/shopping-service/controllers"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

func routeNonAuth() http.Handler {
	r := gin.Default()

	// Set allow CORS
	r.Use(cors.Default())

	// Routes (health check)
	r.GET("/", controllers.HealthCheck)

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

	// Buyer route
	r.POST("/cart", controllers.AddProductToCart)
	r.GET("/cart", controllers.GetCartItems)
	r.PATCH("/cart", controllers.UpdateCartItem)
	r.DELETE("/cart", controllers.DropCart)
	r.GET("/cart/checkout", controllers.Checkout)

	return r
}

func main() {
	// Connect database
	db := config.ConnectDatabase()
	databaseSQL, _ := db.DB()
	defer databaseSQL.Close()

	serverNonAuth := &http.Server{
		Addr:    ":8080",
		Handler: routeNonAuth(),
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
