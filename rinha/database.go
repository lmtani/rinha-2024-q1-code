package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initializeDatabase(ctx context.Context, dbpool *pgxpool.Pool) *pgxpool.Pool {
	var err error
	for i := 0; i < 5; i++ { // Retry up to 5 times
		dbpool, err = pgxpool.New(ctx, os.Getenv("DB_HOSTNAME"))
		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database after retries: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}
