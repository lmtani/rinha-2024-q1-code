package services

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmtani/rinha-2024-q1-code/internal/models"
	"github.com/lmtani/rinha-2024-q1-code/internal/repositories"
)

type Service struct {
	dbpool *pgxpool.Pool
}

var (
	ErrInvalidTransactionType     = errors.New("invalid transaction type")
	ErrClientNotFound             = errors.New("client not found")
	ErrInvalidBalance             = errors.New("invalid balance")
	ErrInsertTransaction          = errors.New("error inserting transaction")
	ErrUpdateSaldo                = errors.New("error updating saldo")
	ErrorInvalidDescription       = errors.New("invalid description")
	ErrorInvalidDescriptionLength = errors.New("invalid description length")
	ErrorInvalidValue             = errors.New("invalid value")
)

func NewService(dbpool *pgxpool.Pool) *Service {
	return &Service{dbpool}
}

func (ts *Service) HandlePostTransactions(clientID int, input models.TransactionInputs) (models.TransactionResponse, error) {
	err := validateInputs(input)
	if err != nil {
		return models.TransactionResponse{}, err
	}

	// Start a new transaction
	tx, err := ts.dbpool.Begin(context.Background())
	if err != nil {
		return models.TransactionResponse{}, err
	}
	defer tx.Rollback(context.Background())

	cliente, err := repositories.GetClient(tx, clientID)
	if err != nil {
		return models.TransactionResponse{}, ErrClientNotFound
	}

	value := newBalance(cliente, input)

	if cliente.Balance+value < -cliente.Limit {
		return models.TransactionResponse{}, ErrInvalidBalance
	}

	err = repositories.InsertTransaction(tx, models.Transaction{ // insertTransaction now uses tx
		ClienteID:   clientID,
		Value:       input.Value,
		Type:        input.Type,
		Description: input.Description,
	})
	if err != nil {
		return models.TransactionResponse{}, ErrInsertTransaction
	}

	// update cliente with the new saldo
	err = repositories.UpdateSaldo(tx, clientID, input.Value)
	err = tx.Commit(context.Background())
	if err != nil {
		return models.TransactionResponse{}, ErrUpdateSaldo
	}

	return models.TransactionResponse{
		Limit:   cliente.Limit,
		Balance: cliente.Balance + value,
	}, nil
}

func validateInputs(input models.TransactionInputs) error {
	if input.Description == "" {
		return ErrorInvalidDescription
	}
	if len(input.Description) > 10 {
		return ErrorInvalidDescriptionLength
	}
	if input.Value <= 0 {
		return ErrorInvalidValue
	}
	if input.Type != "c" && input.Type != "d" {
		return ErrInvalidTransactionType
	}
	return nil
}

func newBalance(client models.Client, t models.TransactionInputs) int {
	var value int
	if t.Type == "c" {
		value = t.Value
	} else {
		value = -t.Value
	}
	return value
}
