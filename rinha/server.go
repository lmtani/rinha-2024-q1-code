package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

var dbpool *pgxpool.Pool

func main() {
	ctx := context.Background()
	dbpool = initializeDatabase(ctx, dbpool)
	defer dbpool.Close()

	router := routing.New()
	router.Get("/clientes/<id>/extrato", handleGetExtrato)
	router.Post("/clientes/<id>/transacoes", handlePostTransacoes)

	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:8080", router.HandleRequest))
}
