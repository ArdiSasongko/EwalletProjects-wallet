package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand/v2"

	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/model"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletService struct {
	q  *sqlc.Queries
	db *pgxpool.Pool
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

func (s *WalletService) WalletCredit(ctx context.Context, payload *model.TransactionCredit) (sqlc.InsertWalletTransactionsCreditRow, error) {
	resp, err := s.createCredit(ctx, int64(payload.UserID), payload.Amount)
	if err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}

	return resp, nil
}

func (s *WalletService) WalletDebit(ctx context.Context, payload *model.TransactionDebit) (sqlc.InsertWalletTransactionsDebitRow, error) {
	resp, err := s.createDebit(ctx, int64(payload.UserID), payload.Amount)
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, err
	}

	return resp, nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID int32) (sqlc.Wallet, error) {
	resp, err := s.q.GetWalletByUserId(ctx, userID)
	if err != nil {
		return sqlc.Wallet{}, err
	}

	return resp, nil
}

func (s *WalletService) GetHistoryTransaction(ctx context.Context, payload model.HistoryPayload) ([]sqlc.GetHistoryTransactionsRow, error) {
	pageSize := payload.Limit
	pageNumber := payload.Offset

	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	allowType := map[string]bool{
		"debit":  true,
		"credit": true,
	}

	if !allowType[payload.TransactionType] {
		payload.TransactionType = "credit"
	}

	resp, err := s.q.GetHistoryTransactions(ctx, sqlc.GetHistoryTransactionsParams{
		WalletID:               payload.WalletID,
		WalletTransactionsType: sqlc.WalletTransactionsType(payload.TransactionType),
		Limit:                  limit,
		Offset:                 offset,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func generateRandomReference(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func (s *WalletService) createCredit(ctx context.Context, userID int64, amount int32) (sqlc.InsertWalletTransactionsCreditRow, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}
	defer tx.Rollback(ctx)

	log.Println(amount)

	qtx := s.q.WithTx(tx)

	amountBigInt := big.NewInt(int64(amount))
	id, err := qtx.CreditWalletBalance(ctx, sqlc.CreditWalletBalanceParams{
		UserID: int32(userID),
		Balance: pgtype.Numeric{
			Int:   amountBigInt,
			Valid: true,
		},
	})
	if err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}

	ref := generateRandomReference(12)
	resp, err := qtx.InsertWalletTransactionsCredit(ctx, sqlc.InsertWalletTransactionsCreditParams{
		WalletID: id,
		Amount: pgtype.Numeric{
			Int:   amountBigInt,
			Valid: true,
		},
		Reference: ref,
	})

	if err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}

	return resp, nil
}

func (s *WalletService) createDebit(ctx context.Context, userID int64, amount int32) (sqlc.InsertWalletTransactionsDebitRow, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)
	amountBigInt := big.NewInt(int64(amount))

	// check balance
	respBal, err := qtx.GetWalletByUserId(ctx, int32(userID))
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, nil
	}

	balance := respBal.Balance.Int.Int64()
	log.Println(balance)
	if balance < amountBigInt.Int64() {
		return sqlc.InsertWalletTransactionsDebitRow{}, fmt.Errorf("balance is not enough for this transaction, balance: %d", balance)
	}

	id, err := qtx.DebitWalletBalance(ctx, sqlc.DebitWalletBalanceParams{
		UserID: int32(userID),
		Balance: pgtype.Numeric{
			Int:   amountBigInt,
			Valid: true,
		},
	})
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, err
	}

	ref := generateRandomReference(12)
	resp, err := qtx.InsertWalletTransactionsDebit(ctx, sqlc.InsertWalletTransactionsDebitParams{
		WalletID: id,
		Amount: pgtype.Numeric{
			Int:   amountBigInt,
			Valid: true,
		},
		Reference: ref,
	})

	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, nil
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, nil
	}

	return resp, nil
}
