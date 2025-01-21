package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type WalletPayload struct {
	UserID int32 `json:"user_id" validate:"required"`
}

func (u *WalletPayload) Validate() error {
	return Validate.Struct(u)
}

type WalletResponse struct {
	UserID    int32     `json:"user_id"`
	Balance   float32   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionCredit struct {
	UserID int32 `json:"user_id"`
	Amount int32 `json:"amount" validate:"required,numeric,gte=50000"`
}

func (u *TransactionCredit) Validate() error {
	return Validate.Struct(u)
}

type TransactionDebit struct {
	UserID int32 `json:"user_id"`
	Amount int32 `json:"amount" validate:"required,numeric,gte=1000"`
}

func (u *TransactionDebit) Validate() error {
	return Validate.Struct(u)
}

type TransactionResponse struct {
	UserID    int32     `json:"user_id"`
	Amount    int32     `json:"amount"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
}

type HistoryPayload struct {
	WalletID        int32  `json:"wallet_id"`
	TransactionType string `json:"transaction_type"`
	Offset          int32  `json:"offset"`
	Limit           int32  `json:"limit"`
}
