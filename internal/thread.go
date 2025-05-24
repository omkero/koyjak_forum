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

type ThreadType struct {
	ThreadID      int64     `json:"thread_id"`
	UserID        int       `json:"user_id"`
	ThreadTitle   string    `json:"thread_title"`
	ThreadContent string    `json:"thread_content"`
	CreatedAt     time.Time `json:"created_at"`
	SafeUrl       string    `json:"safe_url"`
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

func (Th *App) get_thread_controller(ctx *fiber.Ctx) error {
	param := ctx.Params("thread")
	subsParam := strings.Replace(param, "-", " ", -1)
	fmt.Println(subsParam)

	thread, err := Th.get_thread_by_title(subsParam)
	if err != nil {
		fmt.Println(err)
		return ctx.Render("thread", fiber.Map{
			"Error": err.Error(),
		})
	}

	member, err := Th.member_global_information(thread.UserID)
	if err != nil {
		fmt.Println(err)
		return ctx.Render("thread", fiber.Map{
			"Error": err.Error(),
		})
	}

	return ctx.Render("thread", fiber.Map{
		"ThreadTitle":   thread.ThreadTitle,
		"ThreadContent": thread.ThreadContent,
		"UserName":      member.UserName,
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
	var sql_query string = "SELECT * FROM Threads WHERE thread_title = $1"

	err := config.Pool.QueryRow(context.Background(), sql_query, thread_title).Scan(&thread.ThreadID, &thread.UserID, &thread.ThreadTitle, &thread.ThreadContent, &thread.CreatedAt, &thread.SafeUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ThreadType{}, fmt.Errorf("thread not found.")

		}
		fmt.Println(err)
		return ThreadType{}, functions.Something_wnt_wrong()
	}

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

func (Th *App) thread_info_details() {

}
