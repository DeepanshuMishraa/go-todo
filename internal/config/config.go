package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL string
	PORT         string
	JWT_SECRET   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	jwt_secret := os.Getenv("JWT_SECRET")

	if err != nil {
		log.Println("Failed to Load the env file.")
	}

	config := &Config{
		DATABASE_URL: dbUrl,
		PORT:         port,
		JWT_SECRET:   jwt_secret,
	}

	return config, nil
}
