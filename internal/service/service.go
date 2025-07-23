package service

import (
	"context"
	"subscription_service/internal/repo"
	"time"
)

type (
	SubscriptionInput struct {
		ServiceName string
		Price       int
		UserId      string
		StartDate   time.Time
		EndDate     *time.Time
	}

	SubscriptionOutput struct {
		Id          int     `json:"id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserId      string  `json:"user_id"`
		StartDate   string  `json:"start_date"`
		EndDate     *string `json:"end_date"`
	}

	PriceInput struct {
		ServiceName string
		UserId      string
		StartDate   time.Time
		EndDate     time.Time
	}
)

type Subscription interface {
	Create(ctx context.Context, input SubscriptionInput) error
	FindById(ctx context.Context, id int) (SubscriptionOutput, error)
	FindAll(ctx context.Context) ([]SubscriptionOutput, error)
	FindPrice(ctx context.Context, input PriceInput) (int, error)
	Update(ctx context.Context, id int, input SubscriptionInput) error
	Delete(ctx context.Context, id int) error
}

type Services struct {
	Subscription Subscription
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		Subscription: newSubscriptionService(d.Repos.Subscription),
	}
}
