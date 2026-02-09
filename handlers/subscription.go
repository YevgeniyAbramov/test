package handlers

import (
	"log"
	"test/models"
	"test/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionHandler struct {
	subscriptionService *services.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

// CreateSubscription создаёт подписку
// @Summary      Создать подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        body  body  models.CreateSubscriptionRequest  true  "Тело запроса"
// @Success      200  {object}  models.SubscriptionResponse  "Успеx"
// @Failure      400  {object}  models.ErrorResponse  "Невалидный запрос"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /api/v1/subscriptions/ [post]
func (h *SubscriptionHandler) CreateSubscription(c *fiber.Ctx) error {
	var request models.CreateSubscriptionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: "invalid request: " + err.Error(),
		})
	}

	subscription := &models.Subscription{
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := subscription.Validate(); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	err := h.subscriptionService.CreateSubscription(subscription)
	if err != nil {
		log.Printf("[ERROR CREATE] User=%s Error=%v", request.UserID, err)
		return c.Status(500).JSON(models.ErrorResponse{
			Status:  false,
			Message: "failed to create subscription: " + err.Error(),
		})
	}

	log.Printf("[CREATE] ID=%d User=%s Service=%s Price=%d",
		subscription.ID, subscription.UserID, subscription.ServiceName, subscription.Price)

	return c.JSON(models.SubscriptionResponse{
		Status:  true,
		Message: "success",
		Data:    *subscription,
	})
}

// GetSubscription возвращает подписку по ID
// @Summary      Получить подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  models.SubscriptionResponse  "Успеx"
// @Failure      400  {object}  models.ErrorResponse  "Невалидный ID"
// @Failure      404  {object}  models.ErrorResponse  "Подписка не найдена"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /api/v1/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: "invalid request: " + err.Error(),
		})
	}

	subscription, err := h.subscriptionService.GetSubscription(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			log.Printf("[GET] Subscription not found: ID=%d", id)
			return c.Status(404).JSON(models.ErrorResponse{
				Status:  false,
				Message: "subscription not found",
			})
		}
		log.Printf("[ERROR GET] ID=%d Error=%v", id, err)
		return c.Status(500).JSON(models.ErrorResponse{
			Status:  false,
			Message: "failed to get subscription: " + err.Error(),
		})
	}

	log.Printf("[GET] ID=%d", id)

	return c.JSON(models.SubscriptionResponse{
		Status:  true,
		Message: "success",
		Data:    subscription,
	})
}

// DeleteSubscription удаляет подписку (soft delete)
// @Summary      Удалить подписку
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  models.SuccessResponse  "Успеx"
// @Failure      400  {object}  models.ErrorResponse  "Невалидный ID"
// @Failure      404  {object}  models.ErrorResponse  "Подписка не найдена"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /api/v1/subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: "invalid request: " + err.Error(),
		})
	}

	if err := h.subscriptionService.DeleteSubscription(id); err != nil {
		if err.Error() == "subscription not found" {
			return c.Status(404).JSON(models.ErrorResponse{
				Status:  false,
				Message: "subscription not found",
			})
		}
		log.Printf("[ERROR DELETE] ID=%d Error=%v", id, err)
		return c.Status(500).JSON(models.ErrorResponse{
			Status:  false,
			Message: "failed to delete subscription: " + err.Error(),
		})
	}

	log.Printf("[DELETE] ID=%d", id)

	return c.JSON(models.SuccessResponse{
		Status:  true,
		Message: "success",
	})
}

// ListSubscriptions возвращает список подписок с пагинацией
// @Summary      Список подписок
// @Tags         subscriptions
// @Produce      json
// @Param        page   query  int  false  "Страница"   default(1)
// @Param        limit  query  int  false  "Лимит"      default(10)
// @Success      200  {object}  models.ListResponse  "Успеx"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /api/v1/subscriptions/list [get]
func (h *SubscriptionHandler) ListSubscriptions(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	data, err := h.subscriptionService.ListSubscriptions(page, limit)
	if err != nil {
		log.Printf("[ERROR LIST] Page=%d Limit=%d Error=%v", page, limit, err)
		return c.Status(500).JSON(models.ErrorResponse{
			Status:  false,
			Message: "failed to list subscriptions: " + err.Error(),
		})
	}

	log.Printf("[LIST] Page=%d Limit=%d Count=%d Total=%d",
		page, limit, len(data.Subscriptions), data.Total)

	return c.JSON(models.ListResponse{
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

// UpdateSubscription обновляет подписку
// @Summary      Обновить подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "ID подписки"
// @Param        body  body  models.UpdateSubscriptionRequest  true  "Поля для обновления"
// @Success      200  {object}  models.SubscriptionResponse  "Успеx"
// @Failure      400  {object}  models.ErrorResponse  "Невалидный запрос"
// @Failure      404  {object}  models.ErrorResponse  "Подписка не найдена"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /api/v1/subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: "invalid request: " + err.Error(),
		})
	}

	var request models.UpdateSubscriptionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: "invalid request: " + err.Error(),
		})
	}

	data, err := h.subscriptionService.UpdateSubscription(id, request)
	if err != nil {
		if err.Error() == "subscription not found" {
			return c.Status(404).JSON(models.ErrorResponse{
				Status:  false,
				Message: "subscription not found",
			})
		}
		log.Printf("[ERROR UPDATE] ID=%d Error=%v", id, err)
		return c.Status(500).JSON(models.ErrorResponse{
			Status:  false,
			Message: "failed to update subscription: " + err.Error(),
		})
	}

	log.Printf("[UPDATE] ID=%d", id)

	return c.JSON(models.SubscriptionResponse{
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

// GetTotalCost возвращает суммарную стоимость подписок за период
// @Summary      Суммарная стоимость за период
// @Tags         subscriptions
// @Produce      json
// @Param        start         query  string  true   "Начало периода (MM-YYYY)"
// @Param        end           query  string  true   "Конец периода (MM-YYYY)"
// @Param        user_id       query  string  false  "Фильтр по UUID пользователя"
// @Param        service_name  query  string  false  "Фильтр по названию подписки"
// @Success      200  {object}  models.TotalResponse  "Успеx"
// @Failure      400  {object}  models.ErrorResponse  "Невалидные параметры"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /api/v1/subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalCost(c *fiber.Ctx) error {
	var request models.TotalCostRequest
	err := c.QueryParser(&request)
	if err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: "invalid request: " + err.Error(),
		})
	}

	if err := request.Validate(); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	total, err := h.subscriptionService.GetTotalCost(&request)
	if err != nil {
		log.Printf("[ERROR TOTAL] Error=%v", err)
		return c.Status(500).JSON(models.ErrorResponse{
			Status:  false,
			Message: "failed to get total cost: " + err.Error(),
		})
	}

	log.Printf("[TOTAL] Period=%s to %s Total=%d",
		request.PeriodStart, request.PeriodEnd, total.Total)

	return c.JSON(models.TotalResponse{
		Status:  true,
		Message: "success",
		Data:    total,
	})
}
