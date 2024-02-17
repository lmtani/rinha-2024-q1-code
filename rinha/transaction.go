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

func handlePostTransactions(c *routing.Context) error {
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

	if cliente.Balance+value < -cliente.Limit {
		return respondWithError(c, "New saldo exceeds the limit", fasthttp.StatusUnprocessableEntity)
	}

	err = insertTransaction(tx, Transaction{ // insertTransaction now uses tx
		ClienteID:   clientID,
		Value:       value,
		Type:        input.Type,
		Description: input.Description,
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

	return respondWithJSON(c, TransactionResponse{
		Limit:   cliente.Limit,
		Balance: cliente.Balance + value,
	})
}

func parseValue(input TransactionInputs) (int, error) {
	var value int
	if input.Type == "c" {
		value = input.Value
	} else if input.Type == "d" {
		value = -input.Value
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

func parseAndValidateInput(body []byte) (TransactionInputs, error) {
	var input TransactionInputs
	err := json.Unmarshal(body, &input)
	if err != nil {
		return TransactionInputs{}, err
	}
	if input.Description == "" {
		return TransactionInputs{}, errors.New("invalid description")
	}
	if len(input.Description) > 10 {
		return TransactionInputs{}, errors.New("invalid description length")
	}
	if input.Value <= 0 {
		return TransactionInputs{}, errors.New("invalid value")
	}
	return input, nil
}
