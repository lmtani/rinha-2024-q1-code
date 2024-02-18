package services

import (
	"time"
)

func (ts *Service) HandleGetStatement(clientID int) (*StatementResponse, error) {
	cwt, err := ts.repository.GetClientWithTransactions(clientID)
	if err != nil {
		return nil, err
	}

	return &StatementResponse{
		Balance: BalanceResponse{
			Total:       cwt.Balance,
			DataExtrato: time.Now(),
			Limite:      cwt.Limit,
		},
		Transactions: cwt.Transactions,
	}, nil
}
