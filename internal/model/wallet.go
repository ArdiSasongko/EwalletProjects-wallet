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
