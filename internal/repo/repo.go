package repo

import (
	"context"
	"subscription_service/internal/model/dbmodel"
	"subscription_service/internal/repo/pgdb"
	"subscription_service/pkg/postgres"
	"time"
)

type Subscription interface {
	Create(ctx context.Context, s dbmodel.Subscription) error
	FindById(ctx context.Context, id int) (dbmodel.Subscription, error)
	FindAll(ctx context.Context) ([]dbmodel.Subscription, error)
	FindPrice(ctx context.Context, service, userId string, start, end time.Time) (int, error)
	Update(ctx context.Context, s dbmodel.Subscription) error
	Delete(ctx context.Context, id int) error
}

type Repositories struct {
	Subscription
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Subscription: pgdb.NewSubscriptionRepo(pg),
	}
}
