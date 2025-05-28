package internal

import (
	"fmt"
	"koyjak/internal/functions"
	"koyjak/internal/keywords"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func (Th *App) VerifySession(ctx *fiber.Ctx) error {
	sessionKeyScrete := os.Getenv("SESSION_TOKEN_KEY")
	if sessionKeyScrete == "" {
		fmt.Println("SESSION_TOKEN_KEY not set")
		// Optionally: return error or panic
	}

	var cook = ctx.Cookies(keywords.COOKIE_NAME)
	fmt.Println("Session cookie:", cook)

	if cook != "" {
		_, err := JwtVerifySignature(cook, sessionKeyScrete)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)

			functions.DeleteSessionToken(ctx)
			return ctx.JSON(fiber.Map{
				"message": "Invalid or expired session key",
				"status":  http.StatusUnauthorized,
			})
		}
	}

	return ctx.Next()
}

func (Th *App) RedirectIsExist(ctx *fiber.Ctx) error {
	sessionKeyScrete := os.Getenv("SESSION_TOKEN_KEY")
	var cook = ctx.Cookies(keywords.COOKIE_NAME)

	if cook == "" {
		return ctx.Next()
	}

	if cook != "" {
		_, err := JwtVerifySignature(cook, sessionKeyScrete)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)

			functions.DeleteSessionToken(ctx)
			return ctx.Next()

		}
	}

	return ctx.Redirect("/")
}
