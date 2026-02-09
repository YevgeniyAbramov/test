package routes

import (
	"test/handlers"

	"github.com/gofiber/fiber/v2"
)

func Use(app *fiber.App, subscriptionHandler *handlers.SubscriptionHandler) {
	api := app.Group("/api/v1/subscriptions")

	//Подписки
	{
		api.Post("/", subscriptionHandler.CreateSubscription)
		api.Get("/total", subscriptionHandler.GetTotalCost)
		api.Get("/list", subscriptionHandler.ListSubscriptions)
		api.Get("/:id", subscriptionHandler.GetSubscription)
		api.Put("/:id", subscriptionHandler.UpdateSubscription)
		api.Delete("/:id", subscriptionHandler.DeleteSubscription)
	}
}
