package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/order-service/config"
	"github.com/tengkuroman/microshop/order-service/controllers"
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

	// Routes (user)
	r.GET("/orders", controllers.GetOrdersDetail)
	r.DELETE("/order/delete/:order_detail_id", controllers.DeleteOrder)
	r.PATCH("/order/payment/:order_detail_id/:payment_provider_id", controllers.SelectPaymentProvider)
	r.PATCH("/order/payment/checkout/:order_detail_id", controllers.PayOrder)

	return r
}

func routeService(key string, value interface{}) http.Handler {
	r := gin.Default()

	// Set allow CORS
	r.Use(cors.Default())

	// Set context
	r.Use(func(c *gin.Context) {
		c.Set(key, value)
	})

	// Routes (service)
	r.POST("/order", controllers.CreateOrder)

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

	serverService := &http.Server{
		Addr:    ":8082",
		Handler: routeService("db", db),
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

	g.Go(func() error {
		err := serverService.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
