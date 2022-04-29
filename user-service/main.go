package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/user-service/config"
	"github.com/tengkuroman/microshop/user-service/controllers"
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
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

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

	// Routes (registered)
	r.PATCH("/change", controllers.ChangeUserDetail)
	r.PATCH("/change/password/", controllers.ChangePassword)
	r.PATCH("/switch/:role", controllers.SwitchUser)

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

	// Routes (API gateway)
	r.POST("/auth/validate", controllers.ValidateUser)

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
