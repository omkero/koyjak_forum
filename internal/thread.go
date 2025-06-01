package internal

import (
	"context"
	"fmt"
	"koyjak/config"
	"koyjak/internal/functions"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type ThreadForm struct {
	ThreadTitle   string `json:"thread_title" bindind:"required"`
	UserID        int    `json:"user_id"`
	ThreadContent string `json:"thread_content" bindind:"required"`
}

type PostForm struct {
	PostTitle   string `json:"post_title" bindind:"required"`
	UserID      int    `json:"user_id"`
	PostContent string `json:"post_content" bindind:"required"`
	ThreadID    int    `json:"thread_id" bindind:"required"`
	ThreadToken string `json:"thread_token" bindind:"required"`
}

type ThreadType struct {
	ThreadID       int64             `json:"thread_id"`
	UserID         int               `json:"user_id"`
	ThreadTitle    string            `json:"thread_title"`
	ThreadContent  string            `json:"thread_content"`
	CreatedAt      time.Time         `json:"created_at"`
	SafeUrl        string            `json:"safe_url"`
	Member         MemberGlobalModel `json:"member" binding:"required"`
	CreatedAtSince string            `json:"created_at_since"`
	ThreadToken    string            `json:"thread_token"`
}

type ResponseThreadType struct {
	ThreadID      int64             `json:"thread_id"`
	UserID        int               `json:"user_id"`
	ThreadTitle   string            `json:"thread_title"`
	ThreadContent string            `json:"thread_content"`
	CreatedAt     string            `json:"created_at"`
	SafeUrl       string            `json:"safe_url"`
	Member        MemberGlobalModel `json:"member" binding:"required"`
}

type ThreadPost struct {
	PostID         int64             `json:"post_id" bindind:"required"`
	ThreadID       int64             `json:"thread_id" bindind:"required"`
	UserID         int               `json:"user_id" bindind:"required"`
	PostTitle      string            `json:"post_title" bindind:"required"`
	PostContent    string            `json:"post_content" bindind:"required"`
	CreatedAt      time.Time         `json:"created_at" binding:"required"`
	Member         MemberGlobalModel `json:"member" binding:"required"`
	CreatedAtSince string            `json:"created_at_since"`
}

type ThreadsResult struct {
	Threads []ThreadType
	Err     error
}

type ThreadResult struct {
	Thread ThreadType
	Err    error
}

type PostsResult struct {
	Posts []ThreadPost
	Err   error
}

func (Th *App) get_thread_controller(ctx *fiber.Ctx) error {
	param := ctx.Params("thread")
	subsParam := strings.Replace(param, "-", " ", -1)

	isAuthChannel := make(chan IsAuthRsult)
	threadChannel := make(chan ThreadResult)
	latestPostsChannel := make(chan PostsResult)
	threadsChannel := make(chan ThreadsResult)

	go func() {
		member, isAuth, err := Th.is_Auth(ctx)
		isAuthChannel <- IsAuthRsult{
			IsAuth: isAuth,
			Err:    err,
			Member: member,
		}
	}()

	go func() {
		thread, err := Th.get_thread_by_title(subsParam)
		threadChannel <- ThreadResult{
			Thread: thread,
			Err:    err,
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
		threads, err := Th.get_all_threads(5)

		threadsChannel <- ThreadsResult{
			Threads: threads,
			Err:     err,
		}
	}()

	isAuthResult := <-isAuthChannel
	ThreadResult := <-threadChannel
	latestPostsResult := <-latestPostsChannel
	threadsRsult := <-threadsChannel

	if ThreadResult.Err != nil {
		return ctx.Render("thread/thread", fiber.Map{
			"Thread":      ThreadResult.Thread,
			"Threads":     threadsRsult.Threads,
			"IsAuth":      isAuthResult.IsAuth,
			"Member":      isAuthResult.Member,
			"ThreadError": ThreadResult.Err.Error(),
			"LatestPosts": latestPostsResult.Posts,
		})
	}

	posts, err := Th.thread_posts(ThreadResult.Thread.ThreadID)
	if err != nil {
		return ctx.Render("thread/thread", fiber.Map{
			"Thread":      ThreadResult.Thread,
			"Threads":     threadsRsult.Threads,
			"IsAuth":      isAuthResult.IsAuth,
			"Member":      isAuthResult.Member,
			"PostsError":  err.Error(),
			"LatestPosts": latestPostsResult.Posts,
		})
	}

	if latestPostsResult.Err != nil {
		return ctx.Render("thread/thread", fiber.Map{
			"Thread":  ThreadResult.Thread,
			"Threads": threadsRsult.Threads,

			"Posts":            posts,
			"IsAuth":           isAuthResult.IsAuth,
			"Member":           isAuthResult.Member,
			"LatestPostsError": latestPostsResult.Err.Error(),
		})
	}

	if threadsRsult.Err != nil {
		return ctx.Render("thread/thread", fiber.Map{
			"Thread":       ThreadResult.Thread,
			"Posts":        posts,
			"IsAuth":       isAuthResult.IsAuth,
			"Member":       isAuthResult.Member,
			"ThreadsError": threadsRsult.Err.Error(),
			"LatestPosts":  latestPostsResult.Posts,
		})
	}

	return ctx.Render("thread/thread", fiber.Map{
		"Thread":      ThreadResult.Thread,
		"Threads":     threadsRsult.Threads,
		"Posts":       posts,
		"IsAuth":      isAuthResult.IsAuth,
		"Member":      isAuthResult.Member,
		"LatestPosts": latestPostsResult.Posts,
	})
}

func (Th *App) post_thread_controller(ctx *fiber.Ctx) error {

	var Body ThreadForm
	err := ctx.BodyParser(&Body)
	if err != nil {
		fmt.Println(err)
	}

	isAuthChan := make(chan IsAuthRsult)
	insertChan := make(chan error)

	go func() {
		member, isAuth, err := Th.is_Auth(ctx)
		isAuthChan <- IsAuthRsult{
			IsAuth: isAuth,
			Err:    err,
			Member: member,
		}
	}()

	isAuthresult := <-isAuthChan
	Body.UserID = isAuthresult.Member.UserID

	if isAuthresult.Err != nil {
		ctx.Status(http.StatusUnauthorized) // after inserting set the status as created
		return ctx.JSON(fiber.Map{
			"messsage": "You must be logged-in to create a thread.",
			"status":   http.StatusUnauthorized,
		})
	}

	go func() {
		_, err := Th.insert_thread(Body)
		insertChan <- err
	}()

	isInsertedResult := <-insertChan

	if isInsertedResult != nil {
		fmt.Println(err)
	}

	ctx.Status(http.StatusCreated) // after inserting set the status as created
	return ctx.JSON(fiber.Map{
		"messsage": "thread created",
		"status":   http.StatusCreated,
	})
}

func (Th *App) post_reply_controller(ctx *fiber.Ctx) error {
	var Body PostForm
	err := ctx.BodyParser(&Body)
	if err != nil {
		fmt.Println(err)
	}

	if Body.PostTitle == "" {
		ctx.Status(http.StatusBadRequest) // after inserting set the status as created
		return ctx.JSON(fiber.Map{
			"messsage": "missing post_title",
			"status":   http.StatusBadRequest,
		})
	}

	if Body.PostContent == "" {
		ctx.Status(http.StatusBadRequest) // after inserting set the status as created
		return ctx.JSON(fiber.Map{
			"messsage": "missing post_content",
			"status":   http.StatusBadRequest,
		})
	}

	if Body.ThreadToken == "" {
		ctx.Status(http.StatusBadRequest) // after inserting set the status as created
		return ctx.JSON(fiber.Map{
			"messsage": "missing thread_token",
			"status":   http.StatusBadRequest,
		})
	}

	isAuthChan := make(chan IsAuthRsult)
	go func() {
		member, isAuth, err := Th.is_Auth(ctx)
		isAuthChan <- IsAuthRsult{
			IsAuth: isAuth,
			Err:    err,
			Member: member,
		}
	}()

	isAuthresult := <-isAuthChan

	if isAuthresult.Err != nil {
		ctx.Status(http.StatusUnauthorized) // after inserting set the status as created
		return ctx.JSON(fiber.Map{
			"messsage": "You must be logged-in to create a post.",
			"status":   http.StatusUnauthorized,
		})
	}

	threadSecreteKEy := os.Getenv("THREAD_KEY_NAME")
	thread_id, err := VerifyThreadTokenSignature(Body.ThreadToken, threadSecreteKEy)
	if err != nil {
		ctx.Status(http.StatusBadRequest) // after inserting set the status as created
		return ctx.JSON(fiber.Map{
			"messsage": err.Error(),
			"status":   http.StatusBadRequest,
		})
	}

	Body.UserID = isAuthresult.Member.UserID
	Body.ThreadID = thread_id

	_, err = Th.insert_post(Body)
	if err != nil {
		fmt.Println(err)
	}

	ctx.Status(http.StatusCreated) // after inserting set the status as created
	return ctx.JSON(fiber.Map{
		"messsage": "post created",
		"status":   http.StatusCreated,
	})
}

// make sure to inner join member to each one
func (Th *App) get_all_threads(limit int) ([]ThreadType, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var threads []ThreadType
	var sql_query string = "SELECT * FROM Threads ORDER BY created_at DESC LIMIT $1"
	row, err := config.Pool.Query(context.Background(), sql_query, limit)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []ThreadType{}, fmt.Errorf("threads not found !!")
		}
		fmt.Println(err)
		return []ThreadType{}, functions.Something_wnt_wrong()
	}

	for row.Next() {
		var tempThread ThreadType

		err = row.Scan(&tempThread.ThreadID, &tempThread.UserID, &tempThread.ThreadTitle, &tempThread.ThreadContent, &tempThread.CreatedAt, &tempThread.SafeUrl, &tempThread.ThreadToken)
		if err != nil {
			fmt.Println(err)
		}

		t := tempThread.CreatedAt
		date := functions.TimeAgo(t)
		tempThread.CreatedAtSince = date

		threads = append(threads, tempThread)
	}

	return threads, err
}

func (Th *App) get_thread_by_id(thread_id int) ThreadType {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var thread ThreadType
	var sql_query string = "SELECT * FROM Threads WHERE thread_id = $1"

	err := config.Pool.QueryRow(context.Background(), sql_query, thread_id).Scan(thread.ThreadID, thread.UserID, thread.ThreadTitle, thread.ThreadContent, thread.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}

	return thread
}
func (Th *App) get_thread_by_title(thread_title string) (ThreadType, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var thread ThreadType
	//	var sql_query string = "SELECT * FROM Threads WHERE thread_title = $1"

	var sql_query = `
	SELECT 
	   t.thread_id, t.user_id, t.thread_title, t.thread_content, t.created_at, t.safe_url, t.thread_token,
	   u.user_id, u.username, u.email_address, u.created_at
	FROM Threads t
	INNER JOIN Users u ON u.user_id = t.user_id WHERE t.thread_title = $1
	`

	err := config.Pool.QueryRow(context.Background(), sql_query, thread_title).Scan(
		&thread.ThreadID, &thread.UserID, &thread.ThreadTitle, &thread.ThreadContent, &thread.CreatedAt, &thread.SafeUrl, &thread.ThreadToken,
		&thread.Member.UserID, &thread.Member.UserName, &thread.Member.EmailAddress, &thread.Member.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ThreadType{}, fmt.Errorf("thread not found.")
		}
		return ThreadType{}, functions.Something_wnt_wrong()
	}

	t := thread.CreatedAt
	date := functions.TimeAgo(t)
	thread.CreatedAtSince = date

	return thread, nil
}

func (Th *App) insert_thread(body ThreadForm) (bool, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	// there is problem here thread not been create fix it
	var thread_id int
	var safe_url string = strings.Replace(body.ThreadTitle, " ", "-", -1)
	var sql_query string = `
	INSERT INTO Threads (thread_title, user_id, thread_content, safe_url) VALUES ($1, $2, $3, $4)
	RETURNING thread_id
	`

	err := config.Pool.QueryRow(context.Background(), sql_query, body.ThreadTitle, body.UserID, body.ThreadContent, safe_url).Scan(&thread_id)
	if err != nil {
		return false, err
	}
	secreteKey := os.Getenv("THREAD_KEY_NAME")
	threadToken, err := GenerateThreadToken(thread_id, secreteKey)
	if err != nil {
		fmt.Println(err)
	}

	var update_sql_query string = `UPDATE Threads SET thread_token = $1 WHERE thread_id = $2`
	Exec, err := config.Pool.Exec(context.Background(), update_sql_query, threadToken, thread_id)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(Exec.RowsAffected())

	return Exec.RowsAffected() >= 1, nil
}

func (Th *App) thread_posts(thread_id int64) ([]ThreadPost, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var posts []ThreadPost
	sql_query := `
	SELECT 
		p.post_id, p.thread_id, p.user_id, p.post_title, p.post_content, p.created_at,
		u.user_id, u.username, u.email_address, u.created_at
	FROM Posts p
	INNER JOIN Users u ON u.user_id = p.user_id
	WHERE p.thread_id = $1
	`

	rows, err := config.Pool.Query(context.Background(), sql_query, thread_id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []ThreadPost{}, fmt.Errorf("post not found.")
		}
		fmt.Println(err)
		return []ThreadPost{}, functions.Something_wnt_wrong()
	}
	defer rows.Close()

	for rows.Next() {
		var post ThreadPost

		err = rows.Scan(
			&post.PostID, &post.ThreadID, &post.UserID, &post.PostTitle, &post.PostContent, &post.CreatedAt,
			&post.Member.UserID, &post.Member.UserName, &post.Member.EmailAddress, &post.Member.CreatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return []ThreadPost{}, functions.Something_wnt_wrong()
		}

		t := post.CreatedAt
		date := functions.TimeAgo(t)
		post.CreatedAtSince = date

		posts = append(posts, post)
	}

	return posts, nil
}

func (Th *App) insert_post(body PostForm) (bool, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var sql_query string = "INSERT INTO Posts (thread_id, user_id, post_title, post_content) VALUES ($1, $2, $3, $4)"
	Exec, err := config.Pool.Exec(context.Background(), sql_query, body.ThreadID, body.UserID, body.PostTitle, body.PostContent)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return Exec.RowsAffected() >= 1, nil
}

func (Th *App) latest_posts() ([]ThreadPost, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var posts []ThreadPost
	sql_query := `
        SELECT 
        	p.post_id, p.thread_id, p.user_id, p.post_title, p.post_content, p.created_at,
        	u.user_id, u.username, u.email_address, u.created_at
        FROM Posts p
        INNER JOIN Users u ON u.user_id = p.user_id
        ORDER BY p.created_at DESC
        LIMIT 5
	`

	rows, err := config.Pool.Query(context.Background(), sql_query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []ThreadPost{}, fmt.Errorf("post not found.")
		}
		fmt.Println(err)
		return []ThreadPost{}, functions.Something_wnt_wrong()
	}
	defer rows.Close()

	for rows.Next() {
		var post ThreadPost

		err = rows.Scan(
			&post.PostID, &post.ThreadID, &post.UserID, &post.PostTitle, &post.PostContent, &post.CreatedAt,
			&post.Member.UserID, &post.Member.UserName, &post.Member.EmailAddress, &post.Member.CreatedAt,
		)
		if err != nil {
			fmt.Println(err)
			return []ThreadPost{}, functions.Something_wnt_wrong()
		}

		t := post.CreatedAt
		date := functions.TimeAgo(t)
		post.CreatedAtSince = date

		posts = append(posts, post)
	}

	return posts, nil
}
