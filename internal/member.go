package internal

import (
	"context"
	"fmt"
	"koyjak/config"
	"koyjak/internal/functions"
	"time"

	"github.com/jackc/pgx/v5"
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

func (Th *App) member_global_information(user_id int) (MemberGlobalModel, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var Member MemberGlobalModel

	const sql_query string = `SELECT user_id, username, email_address, created_at FROM Users WHERE user_id = $1`
	err := config.Pool.QueryRow(context.Background(), sql_query, user_id).Scan(&Member.UserID, &Member.UserName, &Member.EmailAddress, &Member.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return MemberGlobalModel{}, fmt.Errorf("User not found !!")
		}
		return MemberGlobalModel{}, functions.Something_wnt_wrong()

	}

	return Member, nil
}

// working on authentication
