package services

import (
	"errors"
	"github.com/lmtani/rinha-2024-q1-code/internal/models"
	"github.com/lmtani/rinha-2024-q1-code/internal/repositories"
)

type Service struct {
	repository repositories.Repository
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

func NewService(r repositories.Repository) *Service {
	return &Service{repository: r}
}

func (ts *Service) HandlePostTransactions(clientID int, input models.TransactionInputs) (*models.TransactionResponse, error) {
	err := validateInputs(input)
	if err != nil {
		return nil, err
	}

	cliente, err := ts.repository.GetClient(clientID)
	if err != nil {
		return nil, err
	}

	value := SetValueSignal(input)
	cliente.Balance += value

	if input.Type == "d" && cliente.Balance < -cliente.Limit {
		return nil, ErrInvalidBalance
	}

	err = ts.repository.InsertTransaction(value, models.Transaction{ // insertTransaction now uses tx
		ClienteID:   clientID,
		Value:       input.Value,
		Type:        input.Type,
		Description: input.Description,
	})
	if err != nil {
		return nil, ErrInsertTransaction
	}

	return &models.TransactionResponse{
		Limit:   cliente.Limit,
		Balance: cliente.Balance,
	}, nil
}

func validateInputs(input models.TransactionInputs) error {
	if input.Description == "" {
		return ErrorInvalidDescription
	}
	if len(input.Description) > 10 || len(input.Description) < 1 {
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

func SetValueSignal(t models.TransactionInputs) int {
	var value int
	if t.Type == "c" {
		value = t.Value
	} else {
		value = -t.Value
	}
	return value
}
