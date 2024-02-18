package models

type Client struct {
	ID      int
	Name    string
	Limit   int
	Balance int
}

type ClientWithTransactions struct {
	Client
	Transactions []Transaction
}
