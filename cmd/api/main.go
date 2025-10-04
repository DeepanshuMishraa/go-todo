package main

import (
	"fmt"
	"log"
	"time"

	"github.com/DeepanshuMishraa/gotodo/config"
	"github.com/DeepanshuMishraa/gotodo/internals/auth"
	"github.com/DeepanshuMishraa/gotodo/internals/database"
	"github.com/DeepanshuMishraa/gotodo/internals/handlers"
	"github.com/DeepanshuMishraa/gotodo/internals/middleware"
	"github.com/DeepanshuMishraa/gotodo/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	db, err := database.NewDBConnection(cfg.Database.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, 24*time.Hour)

	userRepo := repository.NewUserRepository(db.DB)
	todoRepo := repository.NewTodoRepository(db.DB)

	authHandler := handlers.NewAuthHandler(userRepo, jwtManager)
	todoHandler := handlers.NewTodoHandler(todoRepo)

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Post("/api/register", authHandler.Register)
	app.Post("/api/login", authHandler.Login)

	api := app.Group("/api", middleware.AuthMiddleware(jwtManager))

	api.Post("/todos", todoHandler.Create)
	api.Get("/todos/:id", todoHandler.GetByID)
	api.Put("/todos/:id", todoHandler.Update)
	api.Delete("/todos/:id", todoHandler.Delete)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	port := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on port %s", cfg.Server.Port)

	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
