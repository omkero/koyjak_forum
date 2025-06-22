package internal

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type App struct{}

func (Th *App) RootPage(ctx *fiber.Ctx) error {
	isAuthChannel := make(chan IsAuthRsult)
	latestPostsChannel := make(chan PostsResult)
	latestThreadsChennel := make(chan ThreadsResult)

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

	go func() {
		latest_threads, err := Th.get_latest_threads(5)

		latestThreadsChennel <- ThreadsResult{
			Threads: latest_threads,
			Err:     err,
		}
	}()

	isAuthResult := <-isAuthChannel
	latestPostsResult := <-latestPostsChannel
	latestThreadsResult := <-latestThreadsChennel

	data, err := Th.get_forums()
	if err != nil {
		return ctx.Render("index", fiber.Map{
			"IsAuth":        isAuthResult.IsAuth,
			"Member":        isAuthResult.Member,
			"LatestPosts":   latestPostsResult.Posts,
			"LatestThreads": latestThreadsResult.Threads,
			"Error":         err.Error(),
		})
	}

	if latestThreadsResult.Err != nil {
		return ctx.Render("index", fiber.Map{
			"IsAuth":             isAuthResult.IsAuth,
			"Member":             isAuthResult.Member,
			"LatestThreadsError": latestThreadsResult.Err.Error(),
			"LatestPosts":        latestPostsResult.Posts,
			"Forums":             data,
		})
	}

	if latestPostsResult.Err != nil {
		return ctx.Render("index", fiber.Map{
			"IsAuth":           isAuthResult.IsAuth,
			"Member":           isAuthResult.Member,
			"LatestPostsError": latestPostsResult.Err.Error(),
			"LatestThreads":    latestThreadsResult.Threads,
			"Forums":           data,
		})
	}

	return ctx.Render("index", fiber.Map{
		"IsAuth":        isAuthResult.IsAuth,
		"Member":        isAuthResult.Member,
		"LatestPosts":   latestPostsResult.Posts,
		"LatestThreads": latestThreadsResult.Threads,
		"Forums":        data,
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

	param := ctx.Params("forumTitle")
	var forum_title string = strings.Replace(param, "-", " ", -1)

	fmt.Println(param)
	fmt.Println(forum_title)

	isAuthChannel := make(chan IsAuthRsult)
	threadsChan := make(chan ThreadsResult)

	go func() {
		member, isAuth, err := Th.is_Auth(ctx)
		isAuthChannel <- IsAuthRsult{
			IsAuth: isAuth,
			Err:    err,
			Member: member,
		}
	}()

	go func() {
		threads, err := Th.get_all_threads_by_forum_title(5, forum_title)
		threadsChan <- ThreadsResult{
			Threads: threads,
			Err:     err,
		}
	}()

	isAuthResult := <-isAuthChannel
	threadsPipe := <-threadsChan
	if isAuthResult.Err != nil {
		fmt.Println(isAuthResult.Err)
	}

	if threadsPipe.Err != nil {
		return ctx.Render("forum/forum", fiber.Map{
			"IsAuth":       isAuthResult.IsAuth,
			"Member":       isAuthResult.Member,
			"ThreadsError": threadsPipe.Err,
		})
	}

	return ctx.Render("forum/forum", fiber.Map{
		"IsAuth":  isAuthResult.IsAuth,
		"Member":  isAuthResult.Member,
		"Threads": threadsPipe.Threads,
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
