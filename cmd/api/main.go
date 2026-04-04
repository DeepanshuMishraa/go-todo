package main

import (
	"log"

	"github.com/DeepanshuMishraa/api-go/internal/config"
	"github.com/DeepanshuMishraa/api-go/internal/database"
	"github.com/DeepanshuMishraa/api-go/internal/handlers"
	"github.com/DeepanshuMishraa/api-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Failed to load the config", err)
	}

	pool, err := database.Connect(cfg.DATABASE_URL)

	if err != nil {
		log.Fatal("Failed to connect to the database", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Welcome to the API",
			"version":  "0.0.1",
			"Configs":  "Loaded",
			"database": "Connected",
		})
	})

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))

	protected.POST("/", handlers.CreateTodoHandler(pool))
	protected.GET("/", handlers.GetAllTodosHandler(pool))
	protected.GET("/:id", handlers.GetTodoByIDHandler(pool))
	protected.PUT("/:id", handlers.UpdateTodoHandler(pool))
	protected.DELETE("/:id", handlers.DeleteTodoHandler(pool))

	router.POST("/user/register", handlers.CreateUserHandler(pool))
	router.POST("/user/login", handlers.LoginHandler(pool, cfg))
	// router.GET("/user/:email", handlers.GetUserByEmailHandler(pool))
	// router.GET("/user/:id", handlers.GetUserByIdHandler(pool))

	router.GET("/test", middleware.AuthMiddleware(cfg), handlers.TestMiddlewareHandler())
	router.Run(":" + cfg.PORT)
}
