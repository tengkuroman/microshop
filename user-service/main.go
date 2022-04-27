package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tengkuroman/microshop/user-service/config"
	"github.com/tengkuroman/microshop/user-service/controllers"
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

	// Routes (public)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Routes (registered)
	r.PATCH("/change", controllers.ChangeUserDetail)
	r.PATCH("/change/password/", controllers.ChangePassword)
	r.PATCH("/switch/:role", controllers.SwitchUser)

	// Routes (API gateway)
	r.POST("/auth/validate", controllers.ValidateUser)

	// Run router
	r.Run()
}
