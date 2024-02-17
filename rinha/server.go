package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var dbpool *pgxpool.Pool

func main() {
	ctx := context.Background()
	dbpool = initializeDatabase(ctx, dbpool)
	defer dbpool.Close()

	fmt.Println("Server running on port 8080")

	router := routing.New()
	router.Get("/clientes/<id>/extrato", handleGetStatement)
	router.Post("/clientes/<id>/transacoes", handlePostTransactions)

	//// Setup pprof handler
	//// Wrap the pprofhandler for compatibility with fasthttp-routing
	//router.Get("/debug/pprof/*", func(c *routing.Context) error {
	//	pprofhandler.PprofHandler(c.RequestCtx)
	//	return nil
	//})

	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8080", router.HandleRequest))
}
