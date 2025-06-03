package internal

import (
	"fmt"

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

func (Th *App) ForumPage(ctx *fiber.Ctx) error {
	isAuthChannel := make(chan IsAuthRsult)

	go func() {
		member, isAuth, err := Th.is_Auth(ctx)
		isAuthChannel <- IsAuthRsult{
			IsAuth: isAuth,
			Err:    err,
			Member: member,
		}
	}()

	isAuthResult := <-isAuthChannel
	if isAuthResult.Err != nil {
		fmt.Println(isAuthResult.Err)
	}
	return ctx.Render("forum/forum", fiber.Map{
		"IsAuth": isAuthResult.IsAuth,
		"Member": isAuthResult.Member,
	})
}

/*
Using Cashing

import ("github.com/patrickmn/go-cache")

var (
	globalCache = cache.New(5*time.Minute, 10*time.Minute) // one-time global init
)

func (Th *App) RootPage(ctx *fiber.Ctx) error {
	threadsChannel := make(chan ThreadsResult)
	isAuthChannel := make(chan IsAuthRsult)
	latestPostsChannel := make(chan PostsResult)

	go func() {
		if cached, found := globalCache.Get("threads"); found {
			threadsChannel <- ThreadsResult{Threads: cached.([]ThreadType), Err: nil}
			return
		}

		threads, err := Th.get_all_threads(5)
		if err == nil {
			globalCache.Set("threads", threads, cache.DefaultExpiration)
		}
		threadsChannel <- ThreadsResult{Threads: threads, Err: err}
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
		if cached, found := globalCache.Get("latest_posts"); found {
			latestPostsChannel <- PostsResult{Posts: cached.([]ThreadPost), Err: nil}
			return
		}

		latestPosts, err := Th.latest_posts()
		if err == nil {
			globalCache.Set("latest_posts", latestPosts, cache.DefaultExpiration)
		}
		latestPostsChannel <- PostsResult{Posts: latestPosts, Err: err}
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

*/
