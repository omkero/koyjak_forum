package internal

import "github.com/gofiber/fiber/v2"

func MainHandler(Koyjak *fiber.App) {
	var appko = App{}

	// main handler
	{
		Koyjak.Get("/", appko.AppHandler)
	}

	// thread handlers
	{
		Koyjak.Post("/create_thread", appko.post_thread_controller)
	}

	// authuntication handlers
	{

		// Koyjak.Get("/auth/login")
	}
}
