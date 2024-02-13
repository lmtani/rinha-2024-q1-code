package main

import "time"

type Cliente struct {
	ID     int
	Nome   string
	Limite int
	Saldo  int
}

type Transacao struct {
	ClienteID   int       `json:"-"` // ignore this field when marshalling
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type ClienteComTransacoes struct {
	Cliente
	Transacoes []Transacao
}

type TransacaoInput struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacaoResponse struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type SaldoResponse struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

type ExtratoResponse struct {
	Saldo      SaldoResponse `json:"saldo"`
	Transacoes []Transacao   `json:"ultimas_transacoes"`
}
