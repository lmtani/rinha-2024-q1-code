package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmtani/rinha-de-backend-2024/internal/models"
	"github.com/lmtani/rinha-de-backend-2024/internal/repositories"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type Service struct {
	dbpool *pgxpool.Pool
}

func NewService(dbpool *pgxpool.Pool) *Service {
	return &Service{dbpool}
}

func (ts *Service) HandlePostTransactions(c *routing.Context) error {
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
	tx, err := ts.dbpool.Begin(context.Background())
	if err != nil {
		return respondWithError(c, "Internal Server Error", fasthttp.StatusInternalServerError)
	}
	defer tx.Rollback(context.Background())

	cliente, err := repositories.GetClient(tx, clientID)
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

	err = repositories.InsertTransaction(tx, models.Transaction{ // insertTransaction now uses tx
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
	err = repositories.UpdateSaldo(tx, clientID, value)

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return respondWithJSON(c, models.TransactionResponse{
		Limit:   cliente.Limit,
		Balance: cliente.Balance + value,
	})
}

func parseValue(input models.TransactionInputs) (int, error) {
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
