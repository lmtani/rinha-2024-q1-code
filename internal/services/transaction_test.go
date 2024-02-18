package services

import (
	"errors"
	"github.com/lmtani/rinha-2024-q1-code/internal/models"
	"github.com/lmtani/rinha-2024-q1-code/internal/repositories"
	"reflect"
	"testing"
)

type MockRepository struct {
	client                 models.Client
	err                    error
	transaction            models.Transaction
	clientWithTransactions models.ClientWithTransactions
}

func (m MockRepository) GetClient(clientID int) (models.Client, error) {
	return m.client, m.err
}

func (m MockRepository) InsertTransaction(t models.Transaction) error {
	return m.err
}

func (m MockRepository) UpdateSaldo(clienteID, valor int) error {
	return m.err
}

func (m MockRepository) GetClientWithTransactions(clienteID int) (models.ClientWithTransactions, error) {
	return m.clientWithTransactions, m.err
}

func TestService_HandlePostTransactions(t *testing.T) {
	type fields struct {
		repository repositories.Repository
	}
	type args struct {
		clientID int
		input    models.TransactionInputs
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *models.TransactionResponse
		wantErrType error
	}{
		{
			name: "Test HandlePostTransactions should return error not found",
			fields: fields{
				repository: MockRepository{
					err: ErrClientNotFound,
				},
			},
			args: args{
				clientID: 1,
				input: models.TransactionInputs{
					Value:       100,
					Type:        "c",
					Description: "Test",
				},
			},
			wantErrType: ErrClientNotFound,
		},
		{
			name: "Test HandlePostTransactions should return error when description larger than 10",
			fields: fields{
				repository: MockRepository{
					client: models.Client{
						ID:      1,
						Name:    "Test",
						Limit:   1000,
						Balance: 1000,
					},
				},
			},
			args: args{
				clientID: 1,
				input: models.TransactionInputs{
					Value:       100,
					Type:        "c",
					Description: "aaaaaaaaaaaaaaaaaaaaa",
				},
			},
			wantErrType: ErrorInvalidDescriptionLength,
		},
		{
			name: "Test HandlePostTransactions with success credit transaction",
			fields: fields{
				repository: MockRepository{
					client: models.Client{
						ID:      1,
						Name:    "Test",
						Limit:   1000,
						Balance: 1000,
					},
				},
			},
			args: args{
				clientID: 1,
				input: models.TransactionInputs{
					Value:       100,
					Type:        "c",
					Description: "Test",
				},
			},
			want: &models.TransactionResponse{
				Limit:   1000,
				Balance: 1100,
			},
		},
		{
			name: "Test HandlePostTransactions with success debit transaction",
			fields: fields{
				repository: MockRepository{
					client: models.Client{
						ID:      1,
						Name:    "Test",
						Limit:   1000,
						Balance: 1000,
					},
				},
			},
			args: args{
				clientID: 1,
				input: models.TransactionInputs{
					Value:       100,
					Type:        "d",
					Description: "Test",
				},
			},
			want: &models.TransactionResponse{
				Limit:   1000,
				Balance: 900,
			},
		},
		{
			name: "Test HandlePostTransactions with invalid balance",
			fields: fields{
				repository: MockRepository{
					client: models.Client{
						ID:      1,
						Name:    "Test",
						Limit:   1,
						Balance: 10,
					},
				},
			},
			args: args{
				clientID: 1,
				input: models.TransactionInputs{
					Value:       100,
					Type:        "d",
					Description: "Test",
				},
			},
			wantErrType: ErrInvalidBalance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Service{
				repository: tt.fields.repository,
			}
			got, err := ts.HandlePostTransactions(tt.args.clientID, tt.args.input)
			if (err != nil) && !errors.Is(err, tt.wantErrType) {
				t.Errorf("HandlePostTransactions() error = %v, wantErr %v", err, tt.wantErrType)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandlePostTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
