package handler

import (
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/service"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Health interface {
		CheckHealth(*fiber.Ctx) error
	}
	Wallet interface {
		CreateWallet(*fiber.Ctx) error
	}
}

func NewHandler(q *sqlc.Queries) Handlers {
	service := service.NewService(q)
	return Handlers{
		Health: &HealthHandler{},
		Wallet: &WalletHandler{
			service: service,
		},
	}
}
