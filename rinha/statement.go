package main

import (
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func handleGetStatement(c *routing.Context) error {
	clientID, err := parseClientID(c.Param("id"))
	if err != nil {
		return respondWithError(c, "Invalid client ID", fasthttp.StatusNotFound)
	}

	cwt, err := getClientWithTransactions(dbpool, clientID)
	if err != nil {
		if err.Error() == "client not found" {
			return respondWithError(c, "Client not found", fasthttp.StatusNotFound)
		}
		return err
	}

	return respondWithJSON(c, StatementResponse{
		Balance: BalanceResponse{
			Total:       cwt.Balance,
			DataExtrato: time.Now(),
			Limite:      cwt.Limit,
		},
		Transactions: cwt.Transacoes,
	})
}
