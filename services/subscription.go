package services

import (
	"test/db"
	"test/models"
)

type SubscriptionService struct {
	db *db.DB
}

func NewSubscriptionService(db *db.DB) *SubscriptionService {
	return &SubscriptionService{db: db}
}

func (s *SubscriptionService) CreateSubscription(subscription *models.Subscription) error {
	err := s.db.CreateSubscription(subscription)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubscriptionService) GetSubscription(id int) (models.Subscription, error) {
	data, err := s.db.GetSubscription(id)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (s *SubscriptionService) DeleteSubscription(id int) error {
	if err := s.db.DeleteSubscription(id); err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) ListSubscriptions(page, limit int) (models.ListSubscriptionsResponse, error) {
	subscriptions, total, err := s.db.ListSubscriptions(page, limit)
	if err != nil {
		return models.ListSubscriptionsResponse{}, err
	}

	return models.ListSubscriptionsResponse{
		Subscriptions: subscriptions,
		Total:         total,
	}, nil
}

func (s *SubscriptionService) UpdateSubscription(id int, updateSubscription models.UpdateSubscriptionRequest) (models.Subscription, error) {
	data, err := s.db.UpdateSubscription(id, &updateSubscription)
	if err != nil {
		return models.Subscription{}, err
	}

	return data, nil
}

func (s *SubscriptionService) GetTotalCost(req *models.TotalCostRequest) (models.TotalCostResponse, error) {
	total, err := s.db.GetTotalCost(req)
	if err != nil {
		return models.TotalCostResponse{}, err
	}

	return models.TotalCostResponse{Total: total}, nil
}
