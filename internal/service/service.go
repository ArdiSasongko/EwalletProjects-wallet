package service

import (
	"context"

	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/model"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
)

type Service struct {
	Wallet interface {
		CreateWallet(context.Context, int32) (*model.WalletResponse, error)
	}
}

func NewService(q *sqlc.Queries) Service {
	return Service{
		Wallet: &WalletService{
			q: q,
		}}
}
