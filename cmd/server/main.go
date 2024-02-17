package main

import (
	"context"
	"fmt"
	"github.com/lmtani/rinha-de-backend-2024/internal/services"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var dbpool *pgxpool.Pool

func main() {
	ctx := context.Background()
	dbpool = initializeDatabase(ctx, dbpool)
	defer dbpool.Close()

	service := services.NewService(dbpool)

	fmt.Println("Server running on port 8080")

	router := routing.New()
	router.Get("/clientes/<id>/extrato", service.HandleGetStatement)
	router.Post("/clientes/<id>/transacoes", service.HandlePostTransactions)

	//// Setup pprof handler
	//// Wrap the pprofhandler for compatibility with fasthttp-routing
	//router.Get("/debug/pprof/*", func(c *routing.Context) error {
	//	pprofhandler.PprofHandler(c.RequestCtx)
	//	return nil
	//})

	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8080", router.HandleRequest))
}

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
