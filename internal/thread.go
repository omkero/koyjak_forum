package internal

import (
	"context"
	"fmt"
	"koyjak/config"
	"koyjak/internal/functions"
	"net/http"
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

func (Th *App) get_thread_controller(ctx *fiber.Ctx) error {
	param := ctx.Params("thread")
	subsParam := strings.Replace(param, "-", " ", -1)

	thread, err := Th.get_thread_by_title(subsParam)
	if err != nil {
		return ctx.Render("thread", fiber.Map{
			"Error": err.Error(),
		})
	}

	posts, err := Th.thread_posts(thread.ThreadID)
	if err != nil {
		fmt.Println(err)
		return ctx.Render("thread", fiber.Map{
			"Error": err.Error(),
		})
	}

	return ctx.Render("thread", fiber.Map{
		"Thread": thread,
		"Posts":  posts,
	})
}

func (Th *App) post_thread_controller(ctx *fiber.Ctx) error {

	var Body ThreadForm
	err := ctx.BodyParser(&Body)
	if err != nil {
		fmt.Println(err)
	}
	Body.UserID = 34
	is_inserted, err := Th.insert_thread(Body)
	if err != nil || !is_inserted {
		fmt.Println(err)
	}
	fmt.Println(is_inserted)

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

	fmt.Println(Body)
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

func (Th *App) get_all_threads() ([]ThreadType, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var threads []ThreadType
	var sql_query string = "SELECT * FROM Threads ORDER BY created_at DESC"
	row, err := config.Pool.Query(context.Background(), sql_query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []ThreadType{}, fmt.Errorf("threads not found !!")
		}
		fmt.Println(err)
		return []ThreadType{}, functions.Something_wnt_wrong()
	}

	for row.Next() {
		var tempThread ThreadType

		err = row.Scan(&tempThread.ThreadID, &tempThread.UserID, &tempThread.ThreadTitle, &tempThread.ThreadContent, &tempThread.CreatedAt, &tempThread.SafeUrl)
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
	   t.thread_id, t.user_id, t.thread_title, t.thread_content, t.created_at, t.safe_url,
	   u.user_id, u.username, u.email_address, u.created_at
	FROM Threads t
	INNER JOIN Users u ON u.user_id = t.user_id WHERE t.thread_title = $1
	`

	err := config.Pool.QueryRow(context.Background(), sql_query, thread_title).Scan(
		&thread.ThreadID, &thread.UserID, &thread.ThreadTitle, &thread.ThreadContent, &thread.CreatedAt, &thread.SafeUrl,
		&thread.Member.UserID, &thread.Member.UserName, &thread.Member.EmailAddress, &thread.Member.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ThreadType{}, fmt.Errorf("thread not found.")
		}
		fmt.Println(err)
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
	var safe_url string = strings.Replace(body.ThreadTitle, " ", "-", -1)
	var sql_query string = "INSERT INTO Threads (thread_title, user_id, thread_content, safe_url) VALUES ($1, $2, $3, $4)"
	Exec, err := config.Pool.Exec(context.Background(), sql_query, body.ThreadTitle, body.UserID, body.ThreadContent, safe_url)
	if err != nil {
		return false, err
	}

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

	body.ThreadID = 141
	body.UserID = 34

	var sql_query string = "INSERT INTO Posts (thread_id, user_id, post_title, post_content) VALUES ($1, $2, $3, $4)"
	Exec, err := config.Pool.Exec(context.Background(), sql_query, body.ThreadID, body.UserID, body.PostTitle, body.PostContent)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return Exec.RowsAffected() >= 1, nil
}
