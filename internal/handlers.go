package internal

import "github.com/gofiber/fiber/v2"

// each post request must call the controller directly and get will return the view
func MainHandler(Koyjak *fiber.App) {
	var appko = App{}
	{
		Koyjak.Get("/", appko.RootPage)
		Koyjak.Get("/forum", appko.ForumPage)

	}
	{
		Koyjak.Post("/create_thread", appko.post_thread_controller)
		Koyjak.Get("/:thread", appko.ThreadPage)
		Koyjak.Post("/create_post", appko.post_reply_controller)
	}
	{
		Koyjak.Post("/auth/session_logout", appko.logout_session_controller)
		Koyjak.Post("/auth/signup", appko.create_member_controller)
		Koyjak.Get("/auth/signup", appko.RedirectIsExist, appko.SignUpPage)
		Koyjak.Post("/auth/signin", appko.signin_member_controller)
		Koyjak.Get("/auth/signin", appko.RedirectIsExist, appko.SignInPage)
	}

}
