package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/external"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/model"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletService struct {
	q        *sqlc.Queries
	db       *pgxpool.Pool
	external external.External
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

	balance, _ := resp.Balance.Float64Value()
	return &model.WalletResponse{
		UserID:    resp.UserID,
		Balance:   float32(balance.Float64),
		CreatedAt: resp.CreatedAt.Time,
	}, err
}

func (s *WalletService) WalletCredit(ctx context.Context, payload *model.TransactionCredit) (model.TransactionResponse, error) {
	resp, err := s.createCredit(ctx, payload.UserID, payload.Amount, payload.Reference)
	amount, _ := resp.Amount.Float64Value()
	if err != nil {
		return model.TransactionResponse{}, err
	}

	mappingResponse := model.TransactionResponse{
		UserID:    resp.WalletID,
		Amount:    amount.Float64,
		Reference: resp.Reference,
		CreatedAt: resp.CreatedAt.Time,
	}

	log.Println(payload.Status)
	templateName := "topup_success"
	if strings.Contains(payload.Reference, "REFUND") {
		templateName = "refund_success"
	}

	log.Println("template", templateName)

	if err := s.external.Notif.SendNotification(ctx, external.NotifRequest{
		Recipient:    payload.Email,
		TemplateName: templateName,
		Placeholder: map[string]string{
			"user_id":    string(resp.WalletID),
			"amount":     fmt.Sprintf("%.2f", amount.Float64),
			"reference":  resp.Reference,
			"created_at": resp.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		},
	}); err != nil {
		return model.TransactionResponse{}, fmt.Errorf("failed send success purchase : %w", err)
	}
	return mappingResponse, nil
}

func (s *WalletService) WalletDebit(ctx context.Context, payload *model.TransactionDebit) (model.TransactionResponse, error) {
	resp, err := s.createDebit(ctx, payload.UserID, payload.Amount, payload.Reference)
	amount, _ := resp.Amount.Float64Value()
	if err != nil {
		return model.TransactionResponse{}, err
	}

	mappingResponse := model.TransactionResponse{
		UserID:    resp.WalletID,
		Amount:    amount.Float64,
		Reference: resp.Reference,
		CreatedAt: resp.CreatedAt.Time,
	}

	log.Println(payload.Status)
	var templateName = "purchase_success"
	log.Println("template", templateName)

	if err := s.external.Notif.SendNotification(ctx, external.NotifRequest{
		Recipient:    payload.Email,
		TemplateName: templateName,
		Placeholder: map[string]string{
			"user_id":    string(resp.WalletID),
			"amount":     fmt.Sprintf("%.2f", amount.Float64),
			"reference":  resp.Reference,
			"created_at": resp.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		},
	}); err != nil {
		return model.TransactionResponse{}, fmt.Errorf("failed send success purchase : %w", err)
	}
	return mappingResponse, nil
}

func (s *WalletService) GetBalance(ctx context.Context, userID int32) (model.BalanceResponse, error) {
	resp, err := s.q.GetWalletByUserId(ctx, userID)
	if err != nil {
		return model.BalanceResponse{}, err
	}

	balance, _ := resp.Balance.Float64Value()
	return model.BalanceResponse{
		UserID:    resp.UserID,
		Balance:   float32(balance.Float64),
		CreatedAt: resp.CreatedAt.Time,
		UpdatedAt: resp.UpdatedAt.Time,
	}, nil
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

func roundToTwoDecimalPlaces(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func (s *WalletService) createCredit(ctx context.Context, userID int32, amount float64, ref string) (sqlc.InsertWalletTransactionsCreditRow, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}
	defer tx.Rollback(ctx)

	log.Println(amount)

	qtx := s.q.WithTx(tx)

	amountFloat := roundToTwoDecimalPlaces(amount)
	amountStr := fmt.Sprintf("%.2f", amountFloat)
	amountNumeric := pgtype.Numeric{}
	if err := amountNumeric.Scan(amountStr); err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, fmt.Errorf("failed to convert amount to numeric :%w", err)
	}
	id, err := qtx.CreditWalletBalance(ctx, sqlc.CreditWalletBalanceParams{
		UserID:  userID,
		Balance: amountNumeric,
	})
	if err != nil {
		return sqlc.InsertWalletTransactionsCreditRow{}, err
	}

	resp, err := qtx.InsertWalletTransactionsCredit(ctx, sqlc.InsertWalletTransactionsCreditParams{
		WalletID:  id,
		Amount:    amountNumeric,
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

func (s *WalletService) createDebit(ctx context.Context, userID int32, amount float64, ref string) (sqlc.InsertWalletTransactionsDebitRow, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)
	amountFloat := roundToTwoDecimalPlaces(amount)
	amountStr := fmt.Sprintf("%.2f", amountFloat)
	amountNumeric := pgtype.Numeric{}
	if err := amountNumeric.Scan(amountStr); err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, fmt.Errorf("failed to convert amount to numeric :%w", err)
	}

	respBal, err := qtx.GetWalletByUserId(ctx, userID)
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, err
	}

	balanceFloat, _ := respBal.Balance.Float64Value()
	log.Println(balanceFloat.Float64)
	log.Println(amountFloat)
	if balanceFloat.Float64 < amountFloat {
		return sqlc.InsertWalletTransactionsDebitRow{}, fmt.Errorf("balance is not enough for this transaction, balance: %.2f", balanceFloat.Float64)
	}

	if balanceFloat.Float64-amountFloat <= 50000.00 {
		return sqlc.InsertWalletTransactionsDebitRow{}, fmt.Errorf("balance is not enough for this transaction, min balance (50000) balance: %.2f", balanceFloat.Float64)
	}

	id, err := qtx.DebitWalletBalance(ctx, sqlc.DebitWalletBalanceParams{
		UserID:  int32(userID),
		Balance: amountNumeric,
	})
	if err != nil {
		return sqlc.InsertWalletTransactionsDebitRow{}, err
	}

	resp, err := qtx.InsertWalletTransactionsDebit(ctx, sqlc.InsertWalletTransactionsDebitParams{
		WalletID:  id,
		Amount:    amountNumeric,
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
