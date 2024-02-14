package main

import (
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func handleGetExtrato(c *routing.Context) error {
	clientID, err := parseClientID(c.Param("id"))
	if err != nil {
		return respondWithError(c, "Invalid client ID", fasthttp.StatusNotFound)
	}

	cwt, err := getClienteWithTransacoes(dbpool, clientID)
	if err != nil {
		if err.Error() == "client not found" {
			return respondWithError(c, "Client not found", fasthttp.StatusNotFound)
		}
		return err
	}

	return respondWithJSON(c, ExtratoResponse{
		Saldo: SaldoResponse{
			Total:       cwt.Saldo,
			DataExtrato: time.Now(),
			Limite:      cwt.Limite,
		},
		Transacoes: cwt.Transacoes,
	})
}
