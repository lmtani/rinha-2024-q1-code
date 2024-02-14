package main

import (
	"context"
	"errors"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func handlePostTransacoes(c *routing.Context) error {
	clientID, err := parseClientID(c.Param("id"))
	if err != nil {
		return respondWithError(c, "Invalid client ID", fasthttp.StatusNotFound)
	}

	// Get the request body
	input, err := parseAndValidateInput(c.Request.Body())
	if err != nil {
		return respondWithError(c, "Invalid request body", fasthttp.StatusUnprocessableEntity)
	}

	// Start a new transaction
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		return respondWithError(c, "Internal Server Error", fasthttp.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	cliente, err := getClient(tx, clientID)
	if err != nil {
		return respondWithError(c, "Client not found", fasthttp.StatusNotFound)
	}

	value, err := parseValue(input)
	if err != nil {
		return respondWithError(c, err.Error(), fasthttp.StatusUnprocessableEntity)
	}

	if cliente.Saldo+value < -cliente.Limite {
		return respondWithError(c, "New saldo exceeds the limit", fasthttp.StatusUnprocessableEntity)
	}

	err = insertTransaction(tx, Transacao{ // insertTransaction now uses tx
		ClienteID: clientID,
		Valor:     value,
		Tipo:      input.Tipo,
		Descricao: input.Descricao,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Message == "New saldo exceeds the limit" {
			c.Error("New saldo exceeds the limit", fasthttp.StatusUnprocessableEntity)
			return nil
		} else {
			// Handle other errors
			c.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			return nil
		}
	}

	// update cliente with the new saldo
	err = updateSaldo(tx, clientID, value)

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return respondWithJSON(c, TransacaoResponse{
		Limite: cliente.Limite,
		Saldo:  cliente.Saldo + value,
	})
}

func parseValue(input TransacaoInput) (int, error) {
	var value int
	if input.Tipo == "c" {
		value = input.Valor
	} else if input.Tipo == "d" {
		value = -input.Valor
	} else {
		return 0, errors.New("invalid transaction type")
	}
	return value, nil
}

func parseClientID(clientIDStr string) (int, error) {
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		return 0, err
	}
	return clientID, nil
}

func parseAndValidateInput(body []byte) (TransacaoInput, error) {
	var input TransacaoInput
	err := json.Unmarshal(body, &input)
	if err != nil {
		return TransacaoInput{}, err
	}
	if input.Descricao == "" {
		return TransacaoInput{}, errors.New("invalid description")
	}
	if len(input.Descricao) > 10 {
		return TransacaoInput{}, errors.New("invalid description length")
	}
	if input.Valor <= 0 {
		return TransacaoInput{}, errors.New("invalid value")
	}
	return input, nil
}
