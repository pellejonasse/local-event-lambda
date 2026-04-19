package main

import (
	"os"

	localeventlambda "local-event-lambda"

	"github.com/gofiber/fiber/v2"
)

type envConfig struct {
	ServiceName string
}

type Handler struct {
	cfg envConfig
	app *fiber.App
}

func newHandler() *Handler {
	h := &Handler{
		cfg: envConfig{
			ServiceName: os.Getenv("SERVICE_NAME"),
		},
		app: fiber.New(),
	}

	h.app.Get("/hello", h.hello)
	h.app.Get("/health", h.health)

	return h
}

func (h *Handler) hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello from " + h.cfg.ServiceName,
	})
}

func (h *Handler) health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func main() {
	h := newHandler()
	localeventlambda.Start(h.app)
}
