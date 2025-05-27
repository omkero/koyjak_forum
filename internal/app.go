package internal

import (
	"github.com/gofiber/fiber/v2"
)

type App struct{}

func (Th *App) RootPage(ctx *fiber.Ctx) error {
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

func (Th *App) SignUpPage(ctx *fiber.Ctx) error {
	return ctx.Render("auth/signup", nil)
}

func (Th *App) SignInPage(ctx *fiber.Ctx) error {
	return ctx.Render("auth/signin", nil)
}

func (Th *App) ThreadPage(ctx *fiber.Ctx) error {
	return Th.get_thread_controller(ctx)
}
