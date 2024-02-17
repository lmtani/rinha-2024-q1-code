package services

import (
	"time"

	"github.com/lmtani/rinha-de-backend-2024/internal/models"
)

type BalanceResponse struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type StatementResponse struct {
	Balance      BalanceResponse      `json:"saldo"`
	Transactions []models.Transaction `json:"ultimas_transacoes"`
}
