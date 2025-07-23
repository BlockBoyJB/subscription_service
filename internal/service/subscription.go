package service

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"subscription_service/internal/model/dbmodel"
	"subscription_service/internal/repo"
	"subscription_service/internal/repo/pgerrs"
	"time"
)

type subscriptionService struct {
	sub repo.Subscription
}

func newSubscriptionService(subscription repo.Subscription) *subscriptionService {
	return &subscriptionService{
		sub: subscription,
	}
}

func (s *subscriptionService) Create(ctx context.Context, input SubscriptionInput) error {
	err := s.sub.Create(ctx, dbmodel.Subscription{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserId:      input.UserId,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	})
	if err != nil {
		log.Err(err).Interface("input", input).Msg("subscription/Create error create subscription in database")
		return err
	}
	log.Info().Interface("input", input).Msg("subscription/Create create new subscription in database")
	return nil
}

func (s *subscriptionService) FindById(ctx context.Context, id int) (SubscriptionOutput, error) {
	sub, err := s.sub.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return SubscriptionOutput{}, ErrSubscriptionNotFound
		}
		log.Err(err).Int("id", id).Msg("subscription/FindById error find subscription in database")
		return SubscriptionOutput{}, err
	}
	output := SubscriptionOutput{
		Id:          id,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserId:      sub.UserId,
		StartDate:   formatDate(sub.StartDate),
	}
	if sub.EndDate != nil {
		output.EndDate = ptr(formatDate(*sub.EndDate))
	}
	return output, nil
}

func (s *subscriptionService) FindAll(ctx context.Context) ([]SubscriptionOutput, error) {
	subscriptions, err := s.sub.FindAll(ctx)
	if err != nil {
		log.Err(err).Msg("subscription/FindAll error find all subscriptions in database")
		return nil, err
	}
	result := make([]SubscriptionOutput, 0, len(subscriptions))
	for _, sub := range subscriptions {
		output := SubscriptionOutput{
			Id:          sub.Id,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserId:      sub.UserId,
			StartDate:   formatDate(sub.StartDate),
		}
		if sub.EndDate != nil {
			output.EndDate = ptr(formatDate(*sub.EndDate))
		}
		result = append(result, output)
	}
	return result, nil
}

func (s *subscriptionService) FindPrice(ctx context.Context, input PriceInput) (int, error) {
	price, err := s.sub.FindPrice(ctx, input.ServiceName, input.UserId, input.StartDate, input.EndDate)
	if err != nil {
		log.Err(err).Interface("input", input).Msg("subscription/FindPrice error find total price in database")
		return 0, err
	}
	return price, nil
}

func (s *subscriptionService) Update(ctx context.Context, id int, input SubscriptionInput) error {
	err := s.sub.Update(ctx, dbmodel.Subscription{
		Id:          id,
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserId:      input.UserId,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	})
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		log.Err(err).Interface("input", input).Msg("subscription/Update error update subscription in database")
		return err
	}
	log.Info().Interface("input", input).Msg("subscription/Update update subscription by id in database")
	return nil
}

func (s *subscriptionService) Delete(ctx context.Context, id int) error {
	if err := s.sub.Delete(ctx, id); err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		log.Err(err).Int("id", id).Msg("subscription/Delete error delete subscription by id in database")
		return err
	}
	log.Info().Int("id", id).Msg("subscription/Delete delete subscription in database")
	return nil
}

func formatDate(t time.Time) string {
	return t.Format("01-2006")
}

func ptr[T any](t T) *T {
	return &t
}
