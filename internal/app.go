package internal

import (
	"github.com/gofiber/fiber/v2"
)

type App struct{}

func (Th *App) AppHandler(ctx *fiber.Ctx) error {
	threads, err := Th.get_all_threads()
	if err != nil {
		return ctx.Render("index", fiber.Map{
			"Error": err.Error(),
		})
	}

	return ctx.Render("index", fiber.Map{
		"Threads": threads, // ✅ fixed key
	})
}
