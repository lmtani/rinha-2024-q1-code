package main

import "time"

type Client struct {
	ID      int
	Name    string
	Limit   int
	Balance int
}

type Transaction struct {
	ClienteID   int       `json:"-"` // ignore this field when marshalling
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em"`
}

type ClientWithTransactions struct {
	Client
	Transacoes []Transaction
}

type TransactionInputs struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type TransactionResponse struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

type BalanceResponse struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type StatementResponse struct {
	Balance      BalanceResponse `json:"saldo"`
	Transactions []Transaction   `json:"ultimas_transacoes"`
}
