package services

import (
	"time"

	"github.com/lmtani/rinha-2024-q1-code/internal/repositories"
)

func (ts *Service) HandleGetStatement(clientID int) (StatementResponse, error) {
	cwt, err := repositories.GetClientWithTransactions(ts.dbpool, clientID)
	if err != nil {
		return StatementResponse{}, err
	}

	return StatementResponse{
		Balance: BalanceResponse{
			Total:       cwt.Balance,
			DataExtrato: time.Now(),
			Limite:      cwt.Limit,
		},
		Transactions: cwt.Transactions,
	}, nil
}
