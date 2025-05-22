package internal

import (
	"context"
	"fmt"
	"koyjak/config"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ThreadForm struct {
	ThreadTitle   string `json:"thread_title" bindind:"required"`
	UserID        int    `json:"user_id"`
	ThreadContent string `json:"thread_content" bindind:"required"`
}

type ThreadType struct {
	ThreadID      int       `json:"thread_id"`
	UserID        int       `json:"user_id"`
	ThreadTitle   string    `json:"thread_title"`
	ThreadContent string    `json:"thread_content"`
	CreatedAt     time.Time `json:"created_at"`
}

/*
func (Th *Thread) get_threads_controller(ctx *fiber.Ctx) error {

	threads := Th.get_all_threads()

	return ctx.Render("index", fiber.Map{
		"Threads": threads, // ✅ fixed key
	})
}
*/

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

func (Th *App) get_all_threads() []ThreadType {
	if config.Pool == nil {
		log.Fatal("Failed to connect db")
	}

	var threads []ThreadType
	var sql_query string = "SELECT * FROM Threads ORDER BY created_at DESC"
	row, err := config.Pool.Query(context.Background(), sql_query)
	if err != nil {
		fmt.Println(err)
	}

	for row.Next() {
		var tempThread ThreadType
		err = row.Scan(&tempThread.ThreadID, &tempThread.UserID, &tempThread.ThreadTitle, &tempThread.ThreadContent, &tempThread.CreatedAt)
		if err != nil {
			fmt.Println(err)
		}
		threads = append(threads, tempThread)
	}

	return threads
}

func (Th *App) insert_thread(body ThreadForm) (bool, error) {
	if config.Pool == nil {
		log.Fatal("Failed to connect db")
	}

	var sql_query string = "INSERT INTO Threads (thread_title, user_id, thread_content) VALUES ($1, $2, $3) "
	Exec, err := config.Pool.Exec(context.Background(), sql_query, body.ThreadTitle, body.UserID, body.ThreadContent)
	if err != nil {
		return false, err
	}

	return Exec.RowsAffected() >= 1, nil
}
