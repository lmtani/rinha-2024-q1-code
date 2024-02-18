package repositories

import (
	"context"
	"database/sql"
	"errors"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lmtani/rinha-2024-q1-code/internal/models"
)

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func GetClient(tx pgx.Tx, clientID int) (models.Client, error) {
	const query = "SELECT id, nome, limite, saldo FROM clientes WHERE id = $1"
	var c models.Client
	err := tx.QueryRow(context.Background(), query, clientID).Scan(&c.ID, &c.Name, &c.Limit, &c.Balance)
	if err != nil {
		return models.Client{}, err
	}
	return c, nil
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func InsertTransaction(tx pgx.Tx, t models.Transaction) error {
	const query = "INSERT INTO transacoes (cliente_id, valor, realizada_em, descricao, tipo) VALUES ($1, $2, now(), $3, $4)"
	_, err := tx.Exec(context.Background(), query, t.ClienteID, t.Value, t.Description, t.Type)
	return err
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func UpdateSaldo(tx pgx.Tx, clienteID, valor int) error {
	const query = "UPDATE clientes SET saldo = saldo + $1 WHERE id = $2"
	_, err := tx.Exec(context.Background(), query, valor, clienteID)
	return err
}

//goland:noinspection SqlNoDataSourceInspection,SqlResolve
func GetClientWithTransactions(dbpool *pgxpool.Pool, clienteID int) (models.ClientWithTransactions, error) {
	const query = `
    SELECT c.id, c.limite, c.saldo, t.valor, t.tipo, t.descricao, t.realizada_em
    FROM clientes c
    LEFT JOIN transacoes t ON c.id = t.cliente_id
    WHERE c.id = $1 
    ORDER BY t.realizada_em DESC
    LIMIT 10`

	rows, err := dbpool.Query(context.Background(), query, clienteID)
	if err != nil {
		return models.ClientWithTransactions{}, err
	}
	defer rows.Close()

	var result models.ClientWithTransactions
	var hasCliente bool
	for rows.Next() {
		var transacao models.Transaction
		var tipo, desc sql.NullString // Use a nullable type for the Tipo
		var realizadaEm sql.NullTime
		var valor sql.NullInt32

		// Scan the row
		if err := rows.Scan(&result.Client.ID, &result.Client.Limit, &result.Client.Balance, &valor, &tipo, &desc, &realizadaEm); err != nil {
			return models.ClientWithTransactions{}, err
		}

		var intVal int
		if valor.Valid {
			intVal = int(math.Abs(float64(valor.Int32)))
		}

		if tipo.Valid { // Check if Tipo is not null
			transacao = models.Transaction{
				ClienteID:   clienteID,
				Description: desc.String,
				Date:        realizadaEm.Time,
				Type:        tipo.String,
				Value:       intVal,
			}
			result.Transactions = append(result.Transactions, transacao)
		}

		if !hasCliente {
			hasCliente = true
		}
	}
	if !hasCliente {
		return models.ClientWithTransactions{}, errors.New("client not found")
	}
	return result, nil
}
