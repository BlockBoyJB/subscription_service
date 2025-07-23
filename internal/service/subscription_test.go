package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"subscription_service/internal/mocks/repomocks"
	"subscription_service/internal/model/dbmodel"
	"subscription_service/internal/repo/pgerrs"
	"testing"
	"time"
)

func TestSubscriptionService_Create(t *testing.T) {
	type args struct {
		ctx   context.Context
		input SubscriptionInput
	}

	type mockBehaviour func(sub *repomocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				},
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Create(a.ctx, dbmodel.Subscription{
					ServiceName: a.input.ServiceName,
					Price:       a.input.Price,
					UserId:      a.input.UserId,
					StartDate:   a.input.StartDate,
					EndDate:     a.input.EndDate,
				}).Return(nil)
			},
			expectErr: nil,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				input: SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				},
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Create(a.ctx, dbmodel.Subscription{
					ServiceName: a.input.ServiceName,
					Price:       a.input.Price,
					UserId:      a.input.UserId,
					StartDate:   a.input.StartDate,
					EndDate:     a.input.EndDate,
				}).Return(errors.New("some error"))
			},
			expectErr: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := repomocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			s := newSubscriptionService(sub)

			err := s.Create(tc.args.ctx, tc.args.input)

			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestSubscriptionService_FindById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type mockBehaviour func(sub *repomocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectOutput  SubscriptionOutput
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(dbmodel.Subscription{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				}, nil)
			},
			expectOutput: SubscriptionOutput{
				Id:          1,
				ServiceName: "Yandex",
				Price:       1000,
				UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
				StartDate:   "01-2025",
				EndDate:     nil,
			},
			expectErr: nil,
		},
		{
			testName: "correct test with end date",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(dbmodel.Subscription{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     ptr(time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)),
				}, nil)
			},
			expectOutput: SubscriptionOutput{
				Id:          1,
				ServiceName: "Yandex",
				Price:       1000,
				UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
				StartDate:   "01-2025",
				EndDate:     ptr("05-2025"),
			},
			expectErr: nil,
		},
		{
			testName: "not found",
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(dbmodel.Subscription{}, pgerrs.ErrNotFound)
			},
			expectErr: ErrSubscriptionNotFound,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(dbmodel.Subscription{}, errors.New("some error"))
			},
			expectErr: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := repomocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			s := newSubscriptionService(sub)

			actual, err := s.FindById(tc.args.ctx, tc.args.id)

			assert.Equal(t, tc.expectOutput, actual)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestSubscriptionService_Update(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    int
		input SubscriptionInput
	}

	type mockBehaviour func(sub *repomocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				id:  1,
				input: SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				},
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Update(a.ctx, dbmodel.Subscription{
					Id:          a.id,
					ServiceName: a.input.ServiceName,
					Price:       a.input.Price,
					UserId:      a.input.UserId,
					StartDate:   a.input.StartDate,
					EndDate:     a.input.EndDate,
				}).Return(nil)
			},
			expectErr: nil,
		},
		{
			testName: "not found",
			args: args{
				ctx: context.Background(),
				id:  2,
				input: SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				},
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Update(a.ctx, dbmodel.Subscription{
					Id:          a.id,
					ServiceName: a.input.ServiceName,
					Price:       a.input.Price,
					UserId:      a.input.UserId,
					StartDate:   a.input.StartDate,
					EndDate:     a.input.EndDate,
				}).Return(pgerrs.ErrNotFound)
			},
			expectErr: ErrSubscriptionNotFound,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				id:  1,
				input: SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				},
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Update(a.ctx, dbmodel.Subscription{
					Id:          a.id,
					ServiceName: a.input.ServiceName,
					Price:       a.input.Price,
					UserId:      a.input.UserId,
					StartDate:   a.input.StartDate,
					EndDate:     a.input.EndDate,
				}).Return(errors.New("some error"))
			},
			expectErr: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := repomocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			s := newSubscriptionService(sub)

			err := s.Update(tc.args.ctx, tc.args.id, tc.args.input)

			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestSubscriptionService_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type mockBehaviour func(sub *repomocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		expectErr     error
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Delete(a.ctx, a.id).Return(nil)
			},
			expectErr: nil,
		},
		{
			testName: "not found",
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Delete(a.ctx, a.id).Return(pgerrs.ErrNotFound)
			},
			expectErr: ErrSubscriptionNotFound,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *repomocks.MockSubscription, a args) {
				sub.EXPECT().Delete(a.ctx, a.id).Return(errors.New("some error"))
			},
			expectErr: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := repomocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			s := newSubscriptionService(sub)

			err := s.Delete(tc.args.ctx, tc.args.id)

			assert.Equal(t, tc.expectErr, err)
		})
	}
}
