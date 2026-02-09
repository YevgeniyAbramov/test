package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name"`
	Price       int        `db:"price" json:"price"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	StartDate   string     `db:"start_date" json:"start_date"`
	EndDate     *string    `db:"end_date" json:"end_date,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"1500"`
	UserID      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   string    `json:"start_date" example:"01-2026"`
	EndDate     *string   `json:"end_date,omitempty" example:"10-2026"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Spotify"`
	Price       *int    `json:"price,omitempty" example:"500"`
	StartDate   *string `json:"start_date,omitempty" example:"02-2026"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2026"`
}

type ListSubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Total         int            `json:"total"`
}

type TotalCostRequest struct {
	PeriodStart string     `query:"start" json:"period_start"`
	PeriodEnd   string     `query:"end" json:"period_end"`
	UserID      *uuid.UUID `query:"user_id" json:"user_id,omitempty"`
	ServiceName *string    `query:"service_name" json:"service_name,omitempty"`
}

type TotalCostResponse struct {
	Total int `json:"total"`
}

func (r *Subscription) Validate() error {
	if r.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if r.ServiceName == "" {
		return errors.New("service_name is required")
	}

	if r.Price < 0 {
		return errors.New("price must be greater than or equal to 0")
	}

	if r.StartDate == "" {
		return errors.New("start_date is required")
	}

	if _, err := time.Parse("01-2006", r.StartDate); err != nil {
		return errors.New("start_date must be in format MM-YYYY")
	}

	if r.EndDate != nil && *r.EndDate != "" {
		if _, err := time.Parse("01-2006", *r.EndDate); err != nil {
			return errors.New("end_date must be in format MM-YYYY")
		}
	}

	return nil
}

func (r *TotalCostRequest) Validate() error {
	if r.PeriodStart == "" {
		return errors.New("start is required")
	}

	if r.PeriodEnd == "" {
		return errors.New("end is required")
	}

	if _, err := time.Parse("01-2006", r.PeriodStart); err != nil {
		return errors.New("start must be in format MM-YYYY")
	}

	if _, err := time.Parse("01-2006", r.PeriodEnd); err != nil {
		return errors.New("end must be in format MM-YYYY")
	}

	return nil
}

type ErrorResponse struct {
	Status  bool   `json:"status" example:"false"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status  bool   `json:"status" example:"true"`
	Message string `json:"message"`
}

type SubscriptionResponse struct {
	Status  bool         `json:"status" example:"true"`
	Message string       `json:"message"`
	Data    Subscription `json:"data"`
}

type ListResponse struct {
	Status  bool                      `json:"status" example:"true"`
	Message string                    `json:"message"`
	Data    ListSubscriptionsResponse `json:"data"`
}

type TotalResponse struct {
	Status  bool              `json:"status" example:"true"`
	Message string            `json:"message"`
	Data    TotalCostResponse `json:"data"`
}
