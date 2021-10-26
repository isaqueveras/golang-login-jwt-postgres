package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/isaqueveras/golang-login-jwt-postgres/database"
	"github.com/isaqueveras/golang-login-jwt-postgres/routes"
)

func main() {
	database.Connect()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if err := app.Listen(":8000"); err != nil {
		log.Panic(err)
		return
	}
}
