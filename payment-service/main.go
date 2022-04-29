package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/payment-service/config"
	"github.com/tengkuroman/microshop/payment-service/controllers"
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

	// Routes (public)
	r.GET("/payment", controllers.GetPaymentProviders)

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

	// Routes (admin)
	r.POST("/payment", controllers.PostPaymentProvider)
	r.PATCH("/payment/:payment_provider_id", controllers.UpdatePaymentProvider)
	r.DELETE("/payment/:payment_provider_id", controllers.DeletePaymentProvider)

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
	r.POST("/payment/process", controllers.ProcessPayment)

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
