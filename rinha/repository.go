package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"math"
)

func getClient(tx pgx.Tx, clientID int) (Cliente, error) {
	const query = "SELECT id, nome, limite, saldo FROM clientes WHERE id = $1"
	var c Cliente
	err := tx.QueryRow(context.Background(), query, clientID).Scan(&c.ID, &c.Nome, &c.Limite, &c.Saldo)
	if err != nil {
		return Cliente{}, err
	}
	return c, nil
}

func insertTransaction(tx pgx.Tx, t Transacao) error {
	const query = "INSERT INTO transacoes (cliente_id, valor, realizada_em, descricao, tipo) VALUES ($1, $2, now(), $3, $4)"
	_, err := tx.Exec(context.Background(), query, t.ClienteID, t.Valor, t.Descricao, t.Tipo)
	return err
}

func updateSaldo(tx pgx.Tx, clienteID, valor int) error {
	const query = "UPDATE clientes SET saldo = saldo + $1 WHERE id = $2"
	_, err := tx.Exec(context.Background(), query, valor, clienteID)
	return err
}

func getClienteWithTransacoes(dbpool *pgxpool.Pool, clienteID int) (ClienteComTransacoes, error) {
	const query = `
    SELECT c.id, c.limite, c.saldo, t.valor, t.tipo, t.descricao, t.realizada_em
    FROM clientes c
    LEFT JOIN transacoes t ON c.id = t.cliente_id
    WHERE c.id = $1 
    ORDER BY t.realizada_em DESC
    LIMIT 10`

	rows, err := dbpool.Query(context.Background(), query, clienteID)
	if err != nil {
		return ClienteComTransacoes{}, err
	}
	defer rows.Close()

	var result ClienteComTransacoes
	var hasCliente bool
	for rows.Next() {
		var transacao Transacao
		var tipo, desc sql.NullString // Use a nullable type for the Tipo
		var realizadaEm sql.NullTime
		var valor sql.NullInt32

		// Scan the row
		if err := rows.Scan(&result.Cliente.ID, &result.Cliente.Limite, &result.Cliente.Saldo, &valor, &tipo, &desc, &realizadaEm); err != nil {
			return ClienteComTransacoes{}, err
		}

		var intVal int
		if valor.Valid {
			intVal = int(math.Abs(float64(valor.Int32)))
		}

		if tipo.Valid { // Check if Tipo is not null
			transacao = Transacao{
				ClienteID:   clienteID,
				Descricao:   desc.String,
				RealizadaEm: realizadaEm.Time,
				Tipo:        tipo.String,
				Valor:       intVal,
			}
			result.Transacoes = append(result.Transacoes, transacao)
		}

		if !hasCliente {
			hasCliente = true
		}
	}
	if !hasCliente {
		return ClienteComTransacoes{}, errors.New("client not found")
	}
	return result, nil
}
