package service

import (
	"context"

	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/model"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Wallet interface {
		CreateWallet(context.Context, int32) (*model.WalletResponse, error)
		WalletCredit(context.Context, *model.TransactionCredit) (model.TransactionResponse, error)
		WalletDebit(context.Context, *model.TransactionDebit) (model.TransactionResponse, error)
		GetBalance(context.Context, int32) (model.BalanceResponse, error)
		GetHistoryTransaction(ctx context.Context, payload model.HistoryPayload) ([]sqlc.GetHistoryTransactionsRow, error)
	}
}

func NewService(q *sqlc.Queries, db *pgxpool.Pool) Service {
	return Service{
		Wallet: &WalletService{
			q:  q,
			db: db,
		}}
}
