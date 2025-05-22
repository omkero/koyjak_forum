package internal

import "github.com/gofiber/fiber/v2"

type App struct{}

func (Th *App) AppHandler(ctx *fiber.Ctx) error {
	threads := Th.get_all_threads()

	return ctx.Render("index", fiber.Map{
		"Threads": threads, // ✅ fixed key
	})
}
