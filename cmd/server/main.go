package main

import (
	"koyjak/config"
	"koyjak/internal"
	"koyjak/internal/functions"
	"log"

	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	engine := html.New("./views", ".html")
	engine.AddFunc("truncate", functions.Truncate)
	engine.AddFunc("truncateFisrt", functions.TruncateFirstLetter)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.InitDB()

	app.Static("/", "./public")
	internal.MainHandler(app)

	// DI

	app.Listen(":8080")
}
