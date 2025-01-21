package handler

import "github.com/gofiber/fiber/v2"

type Handlers struct {
	Health interface {
		CheckHealth(*fiber.Ctx) error
	}
}

func NewHandler() Handlers {
	return Handlers{
		Health: &HealthHandler{},
	}
}
