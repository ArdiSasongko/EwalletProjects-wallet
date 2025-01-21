package handler

import (
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/external"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/service"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/storage/sqlc"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Health interface {
		CheckHealth(*fiber.Ctx) error
	}
	Wallet interface {
		CreateWallet(*fiber.Ctx) error
		Credit(*fiber.Ctx) error
		Debit(*fiber.Ctx) error
		Balance(*fiber.Ctx) error
		HistoryTransaction(ctx *fiber.Ctx) error
	}
	Middleware interface {
		AuthMiddleware() fiber.Handler
	}
}

func NewHandler(q *sqlc.Queries, db *pgxpool.Pool) Handlers {
	service := service.NewService(q, db)
	userManagement := external.NewUserManagement()
	return Handlers{
		Health: &HealthHandler{},
		Wallet: &WalletHandler{
			service: service,
		},
		Middleware: &MiddlewareHandler{
			userManagement: userManagement,
		},
	}
}
