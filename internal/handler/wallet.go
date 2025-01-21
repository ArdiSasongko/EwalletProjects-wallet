package handler

import (
	"fmt"

	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/config/logger"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/model"
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/service"
	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	service service.Service
}

var log = logger.NewLogger()

func (h *WalletHandler) CreateWallet(ctx *fiber.Ctx) error {
	payload := new(model.WalletPayload)

	if err := ctx.BodyParser(payload); err != nil {
		log.WithError(err).Errorf("bad request error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := payload.Validate(); err != nil {
		errorValidate := fmt.Errorf("validate error")
		log.WithError(errorValidate).Errorf("bad request error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp, err := h.service.Wallet.CreateWallet(ctx.Context(), payload.UserID)
	if err != nil {
		log.WithError(err).Errorf("internal server error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
