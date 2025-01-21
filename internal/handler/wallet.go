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

func (h *WalletHandler) Credit(ctx *fiber.Ctx) error {
	data := ctx.Locals("token").(model.TokenResponse)
	payload := new(model.TransactionCredit)

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

	payload.UserID = data.UserID

	resp, err := h.service.Wallet.WalletCredit(ctx.Context(), payload)
	if err != nil {
		log.WithError(err).Errorf("internal server error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *WalletHandler) Debit(ctx *fiber.Ctx) error {
	data := ctx.Locals("token").(model.TokenResponse)
	payload := new(model.TransactionDebit)

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

	payload.UserID = data.UserID

	resp, err := h.service.Wallet.WalletDebit(ctx.Context(), payload)
	if err != nil {
		log.WithError(err).Errorf("internal server error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *WalletHandler) Balance(ctx *fiber.Ctx) error {
	data := ctx.Locals("token").(model.TokenResponse)

	resp, err := h.service.Wallet.GetBalance(ctx.Context(), data.UserID)
	if err != nil {
		log.WithError(err).Errorf("internal server error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (h *WalletHandler) HistoryTransaction(ctx *fiber.Ctx) error {
	data := ctx.Locals("token").(model.TokenResponse)

	limit := ctx.QueryInt("limit", 5)
	offset := ctx.QueryInt("offset", 1)
	typeTransaction := ctx.Query("type", "credit")

	payload := model.HistoryPayload{
		WalletID:        data.UserID,
		Limit:           int32(limit),
		Offset:          int32(offset),
		TransactionType: typeTransaction,
	}

	resp, err := h.service.Wallet.GetHistoryTransaction(ctx.Context(), payload)
	if err != nil {
		log.WithError(err).Errorf("internal server error, method: %v, path: %v", ctx.Method(), ctx.Path())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
