package internal

import (
	"github.com/gofiber/fiber/v2"
)

type App struct{}

func (Th *App) RootPage(ctx *fiber.Ctx) error {
	threadsChannel := make(chan ThreadsResult)
	isAuthChannel := make(chan IsAuthRsult)
	latestPostsChannel := make(chan PostsResult)

	go func() {
		threads, err := Th.get_all_threads(5)

		threadsChannel <- ThreadsResult{
			Threads: threads,
			Err:     err,
		}
	}()

	go func() {
		member, isAuth, err := Th.is_Auth(ctx)
		isAuthChannel <- IsAuthRsult{
			IsAuth: isAuth,
			Err:    err,
			Member: member,
		}
	}()

	go func() {
		latest_posts, err := Th.latest_posts()

		latestPostsChannel <- PostsResult{
			Posts: latest_posts,
			Err:   err,
		}
	}()

	isAuthResult := <-isAuthChannel
	threadsRsult := <-threadsChannel
	latestPostsResult := <-latestPostsChannel

	if threadsRsult.Err != nil {
		return ctx.Render("index", fiber.Map{
			"Error": threadsRsult.Err.Error(),
		})
	}

	if latestPostsResult.Err != nil {
		return ctx.Render("index", fiber.Map{
			"Threads":          threadsRsult.Threads,
			"IsAuth":           isAuthResult.IsAuth,
			"Member":           isAuthResult.Member,
			"LatestPostsError": latestPostsResult.Err.Error(),
		})
	}

	return ctx.Render("index", fiber.Map{
		"Threads":     threadsRsult.Threads,
		"IsAuth":      isAuthResult.IsAuth,
		"Member":      isAuthResult.Member,
		"LatestPosts": latestPostsResult.Posts,
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
