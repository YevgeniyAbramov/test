//go:generate swag init -g cmd/main.go --parseDependency -o ../docs
package main

import (
	"log"
	"os"
	"test/db"
	_ "test/docs"
	"test/handlers"
	"test/routes"
	"test/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title        Subscriptions API
// @version      1.0
// @description REST API для управления и агрегации пользовательских онлайн-подписок.
// @description Позволяет создавать, обновлять, удалять и получать информацию о подписках.
// @host         localhost:4001
// @BasePath     /
func main() {
	db, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	port := os.Getenv("port")
	if port == "" {
		port = "4001"
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			return ctx.Status(code).JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		},
		DisableStartupMessage: false,
	})

	app.Use(logger.New())

	subscriptionService := services.NewSubscriptionService(db)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

	routes.Use(app, subscriptionHandler)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	log.Fatal(app.Listen(":" + port))

}
