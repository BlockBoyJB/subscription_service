package pgdb

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"subscription_service/internal/model/dbmodel"
	"subscription_service/internal/repo/pgerrs"
	"subscription_service/pkg/postgres"
	"time"
)

const (
	subscriptionTable = "subscription"
)

type SubscriptionRepo struct {
	*postgres.Postgres
}

func NewSubscriptionRepo(pg *postgres.Postgres) *SubscriptionRepo {
	return &SubscriptionRepo{pg}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s dbmodel.Subscription) error {
	sql, args, _ := r.Builder.
		Insert(subscriptionTable).
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(s.ServiceName, s.Price, s.UserId, s.StartDate, s.EndDate).
		ToSql()

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		return err
	}
	return nil
}

func (r *SubscriptionRepo) FindById(ctx context.Context, id int) (dbmodel.Subscription, error) {
	sql, args, _ := r.Builder.
		Select("service_name", "price", "user_id", "start_date", "end_date").
		From(subscriptionTable).
		Where("id = ?", id).
		ToSql()

	var s dbmodel.Subscription

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&s.ServiceName,
		&s.Price,
		&s.UserId,
		&s.StartDate,
		&s.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dbmodel.Subscription{}, pgerrs.ErrNotFound
		}
	}
	return s, nil
}

func (r *SubscriptionRepo) FindAll(ctx context.Context) ([]dbmodel.Subscription, error) {
	sql, args, _ := r.Builder.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From(subscriptionTable).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dbmodel.Subscription

	for rows.Next() {
		var s dbmodel.Subscription

		err = rows.Scan(
			&s.Id,
			&s.ServiceName,
			&s.Price,
			&s.UserId,
			&s.StartDate,
			&s.EndDate,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *SubscriptionRepo) FindPrice(ctx context.Context, service, userId string, start, end time.Time) (int, error) {
	b := r.Builder.
		Select("COALESCE(SUM(price), 0)").
		From(subscriptionTable)

	if userId != "" {
		b = b.Where("user_id = ?", userId)
	}
	if service != "" {
		b = b.Where("service_name = ?", service)
	}
	sql, args, _ := b.Where(squirrel.Expr("start_date <= ? AND (end_date IS NULL OR end_date >= ?)", end, start)).ToSql()

	var price int

	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&price); err != nil {
		return 0, err
	}
	return price, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, s dbmodel.Subscription) error {
	sql, args, _ := r.Builder.
		Update(subscriptionTable).
		Set("service_name", s.ServiceName).
		Set("price", s.Price).
		Set("user_id", s.UserId).
		Set("start_date", s.StartDate).
		Set("end_date", s.EndDate).
		Where("id = ?", s.Id).
		ToSql()

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgerrs.ErrNotFound
	}
	return nil
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id int) error {
	sql, args, _ := r.Builder.
		Delete(subscriptionTable).
		Where("id = ?", id).
		ToSql()

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgerrs.ErrNotFound
	}
	return nil
}
