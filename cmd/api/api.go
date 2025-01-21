package api

import (
	"github.com/ArdiSasongko/EwalletProjects-wallet/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type application struct {
	handler handler.Handlers
	config  Config
}

type Config struct {
	addrHTTP string
	addrGRPC string
	logger   *logrus.Logger
	db       DBConfig
	auth     AuthConfig
}

type DBConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type AuthConfig struct {
	secret string
	iss    string
	aud    string
}

func (app *application) mount() *fiber.App {
	r := fiber.New()

	r.Get("/health", app.handler.Health.CheckHealth)

	return r
}

func (app *application) run(r *fiber.App) error {
	app.config.logger.Printf("http server has running, port%v", app.config.addrHTTP)
	return r.Listen(app.config.addrHTTP)
}
