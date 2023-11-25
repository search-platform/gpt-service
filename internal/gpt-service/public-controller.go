package gptservice

import (
	"github.com/gofiber/fiber/v2"
	api "github.com/search-platform/gpt-service/api/gpt"
)

type PublicController struct {
	gptService api.GptServiceServer
}

func NewPublicController(gptService api.GptServiceServer) *PublicController {
	c := &PublicController{
		gptService: gptService,
	}
	return c
}

func (s *PublicController) RegisterController(app *fiber.App) {
	app.Get("/health")
	gpt := app.Group("gpt")
	gpt.Get("/", s.Health)
}

func (s *PublicController) Health(ctx *fiber.Ctx) error {
	return nil
}
