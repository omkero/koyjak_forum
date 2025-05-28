package functions

import (
	"koyjak/internal/keywords"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DeleteSessionToken(ctx *fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     keywords.COOKIE_NAME,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Expires:  time.Unix(0, 0),
	})
}
