package internal

import "github.com/gofiber/fiber/v2"

func MainHandler(Koyjak *fiber.App) {
	var appko = App{}

	{
		Koyjak.Get("/", appko.AppHandler)
	}
	{
		Koyjak.Post("/create_thread", appko.post_thread_controller)
		Koyjak.Get("/:thread", appko.get_thread_controller)
	}

}
