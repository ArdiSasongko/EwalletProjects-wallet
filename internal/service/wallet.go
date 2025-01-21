package service

import (
	"context"
	"errors"

	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/model"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
)

type WalletService struct {
	q *sqlc.Queries
}

func (s *WalletService) CreateWallet(ctx context.Context, id int32) (*model.WalletResponse, error) {
	resp, err := s.q.CreateWallet(ctx, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, errors.New("resource already exists")
			}
		}
	}

	return &model.WalletResponse{
		UserID:    resp.UserID,
		Balance:   float32(resp.Balance.Exp),
		CreatedAt: resp.CreatedAt.Time,
	}, err
}
