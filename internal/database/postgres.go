package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseUrl string) (*pgxpool.Pool, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(databaseUrl)

	if err != nil {
		log.Println("Unable to parse the db url")
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Connected To Database")
	return pool, nil
}
