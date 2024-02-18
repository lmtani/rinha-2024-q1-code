package main

import (
	"errors"
	"github.com/goccy/go-json"
	"github.com/lmtani/rinha-2024-q1-code/internal/models"
	"github.com/lmtani/rinha-2024-q1-code/internal/repositories"
	"github.com/lmtani/rinha-2024-q1-code/internal/services"
	"github.com/valyala/fasthttp"
	"strconv"
)

var errorResponseMap = map[error]struct {
	Message    string
	StatusCode int
}{
	services.ErrInvalidTransactionType:     {"Invalid transaction type", fasthttp.StatusUnprocessableEntity},
	services.ErrorInvalidDescriptionLength: {"Invalid description length", fasthttp.StatusUnprocessableEntity},
	services.ErrorInvalidDescription:       {"Invalid description", fasthttp.StatusUnprocessableEntity},
	services.ErrorInvalidValue:             {"Invalid value", fasthttp.StatusUnprocessableEntity},
	services.ErrClientNotFound:             {"Client not found", fasthttp.StatusNotFound},
	services.ErrInvalidBalance:             {"Invalid balance", fasthttp.StatusUnprocessableEntity},
	services.ErrInsertTransaction:          {"Error inserting transaction", fasthttp.StatusInternalServerError},
	services.ErrUpdateSaldo:                {"Error updating saldo", fasthttp.StatusInternalServerError},
	repositories.ErrClientNotFound:         {"Client not found", fasthttp.StatusNotFound},
}

func parseAndValidateInput(body []byte) (models.TransactionInputs, error) {
	var input models.TransactionInputs
	err := json.Unmarshal(body, &input)
	if err != nil {
		return models.TransactionInputs{}, err
	}
	if input.Description == "" {
		return models.TransactionInputs{}, errors.New("invalid description")
	}
	if len(input.Description) > 10 {
		return models.TransactionInputs{}, errors.New("invalid description length")
	}
	if input.Value <= 0 {
		return models.TransactionInputs{}, errors.New("invalid value")
	}
	return input, nil
}

func parseClientID(clientIDStr string) (int, error) {
	clientID, err := strconv.Atoi(clientIDStr)
	if err != nil {
		return 0, err
	}
	return clientID, nil
}
