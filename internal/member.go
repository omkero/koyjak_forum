package internal

import (
	"context"
	"errors"
	"fmt"
	"koyjak/config"
	"koyjak/internal/functions"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MemberModel struct {
	UserID       int       `json:"user_id" binding:"required"`
	UserName     string    `json:"username" binding:"required"`
	EmailAddress string    `json:"email_address" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	CreatedAt    time.Time `json:"created_at" bindind:"required"`
}
type MemberGlobalModel struct {
	UserID       int       `json:"user_id" binding:"required"`
	UserName     string    `json:"username" binding:"required"`
	EmailAddress string    `json:"email_address" binding:"required"`
	CreatedAt    time.Time `json:"created_at" bindind:"required"`
}

type MemberBody struct {
	UserName     string `json:"username" binding:"required"`
	EmailAddress string `json:"email_address" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

func (Th *App) create_member_controller(ctx *fiber.Ctx) error {
	var Body MemberBody

	err := ctx.BodyParser(&Body)
	if err != nil {
		fmt.Println(err)
	}

	if Body.UserName == "" {
		ctx.Status(http.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "missing username",
			"status":  http.StatusBadRequest,
		})
	}

	if Body.EmailAddress == "" {
		ctx.Status(http.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "missing email_address",
			"status":  http.StatusBadRequest,
		})
	}

	if Body.Password == "" {
		ctx.Status(http.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "missing password",
			"status":  http.StatusBadRequest,
		})
	}

	if len(Body.Password) < 8 {
		ctx.Status(http.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": "password too short it must be 8 charecters or longer",
			"status":  http.StatusBadRequest,
		})
	}

	err = Th.insert_member(Body)
	if err != nil {
		fmt.Println(err)
		ctx.Status(http.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}

	ctx.Status(http.StatusCreated)

	/*
		Here we set Session Cookies back to the client
	*/

	return ctx.JSON(fiber.Map{
		"message": "member created",
		"status":  http.StatusCreated,
	})
}

func (Th *App) member_global_information(user_id int) (MemberGlobalModel, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var Member MemberGlobalModel

	const sql_query string = `SELECT user_id, username, email_address, created_at FROM Users WHERE user_id = $1`
	err := config.Pool.QueryRow(context.Background(), sql_query, user_id).Scan(&Member.UserID, &Member.UserName, &Member.EmailAddress, &Member.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return MemberGlobalModel{}, fmt.Errorf("user not found")
		}
		return MemberGlobalModel{}, functions.Something_wnt_wrong()

	}

	return Member, nil
}

// working on authentication

func (Th *App) insert_member(Body MemberBody) error {
	if config.Pool == nil {
		log.Fatal("cannot establish connection")
	}

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var sql_query string = `INSERT INTO Users (username, email_address, pwd) VALUES ($1, $2, $3)`
	_, err = config.Pool.Exec(context.Background(), sql_query, Body.UserName, Body.EmailAddress, hashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == `23505` {
				return fmt.Errorf("this user already exist")
			}
		}
		return fmt.Errorf("ops something went wrong")
	}

	return nil
}
