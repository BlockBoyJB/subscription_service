package v1

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"subscription_service/internal/mocks/servicemocks"
	"subscription_service/internal/service"
	"subscription_service/pkg/validator"
	"testing"
	"time"
)

func TestSubscriptionRouter_create(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.SubscriptionInput
	}

	type mockBehaviour func(sub *servicemocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		inputBody     string
		expectCode    int
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: service.SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     ptr(time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().Create(a.ctx, a.input).Return(nil)
			},
			inputBody:  `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "07-2025", "end_date": "10-2025"}`,
			expectCode: http.StatusOK,
		},
		{
			testName: "correct test without end date",
			args: args{
				ctx: context.Background(),
				input: service.SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     nil,
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().Create(a.ctx, a.input).Return(nil)
			},
			inputBody:  `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "07-2025"}`,
			expectCode: http.StatusOK,
		},
		{
			testName:      "missing service field",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "07-2025"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "missing price field",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"service_name": "Yandex", "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "07-2025"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "missing user id field",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"service_name": "Yandex", "price": 1000, "start_date": "07-2025"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "missing start date field",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "user id field is not uuid",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"service_name": "Yandex", "price": 1000, "user_id": "foobar", "start_date": "07-2025"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "invalid start date input",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "2025-07-01"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "invalid end date input",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			inputBody:     `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "07-2025", "end_date": "2025-08-01"}`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				input: service.SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     ptr(time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().Create(a.ctx, a.input).Return(errors.New("some error"))
			},
			inputBody:  `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "07-2025", "end_date": "10-2025"}`,
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := servicemocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			e := echo.New()
			e.Validator = validator.NewValidator()
			NewRouter(e, &service.Services{Subscription: sub})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/subscription", bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
		})
	}
}

func TestSubscriptionRouter_findById(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type mockBehaviour func(sub *servicemocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		inputId       int
		expectBody    string
		expectCode    int
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(service.SubscriptionOutput{
					Id:          a.id,
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   "07-2025",
					EndDate:     nil,
				}, nil)
			},
			inputId:    1,
			expectBody: `{"id":1,"service_name":"Yandex","price":1000,"user_id":"6114696a-d069-4fad-a3ed-f27c13651c3a","start_date":"07-2025","end_date":null}` + "\n",
			expectCode: http.StatusOK,
		},
		{
			testName: "not found",
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(service.SubscriptionOutput{}, service.ErrSubscriptionNotFound)
			},
			inputId:    2,
			expectCode: http.StatusNotFound,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().FindById(a.ctx, a.id).Return(service.SubscriptionOutput{}, errors.New("some error"))
			},
			inputId:    1,
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := servicemocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			e := echo.New()
			e.Validator = validator.NewValidator()
			NewRouter(e, &service.Services{Subscription: sub})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/subscription/"+strconv.Itoa(tc.inputId), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Equal(t, tc.expectBody, w.Body.String())
		})
	}
}

func TestSubscriptionRouter_findPrice(t *testing.T) {
	type args struct {
		ctx   context.Context
		input service.PriceInput
	}

	type mockBehaviour func(sub *servicemocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		query         string
		expectBody    string
		expectCode    int
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				input: service.PriceInput{
					ServiceName: "Yandex",
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().FindPrice(a.ctx, a.input).Return(1000, nil)
			},
			query:      `service_name=Yandex&user_id=6114696a-d069-4fad-a3ed-f27c13651c3a&start=01-2025&end=03-2025`,
			expectBody: `{"price":1000}` + "\n",
			expectCode: http.StatusOK,
		},
		{
			testName:      "incorrect start interval",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			query:         `start=2025-01-01&end=03-2025`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName:      "incorrect end interval",
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {},
			query:         `start=01-2025&end=2025-05-01`,
			expectCode:    http.StatusBadRequest,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				input: service.PriceInput{
					ServiceName: "Yandex",
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().FindPrice(a.ctx, a.input).Return(0, errors.New("some error"))
			},
			query:      `service_name=Yandex&user_id=6114696a-d069-4fad-a3ed-f27c13651c3a&start=01-2025&end=03-2025`,
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := servicemocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			e := echo.New()
			e.Validator = validator.NewValidator()
			NewRouter(e, &service.Services{Subscription: sub})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/subscription/price?"+tc.query, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
			assert.Equal(t, tc.expectBody, w.Body.String())
		})
	}
}

func TestSubscriptionRouter_update(t *testing.T) {
	type args struct {
		ctx   context.Context
		id    int
		input service.SubscriptionInput
	}

	type mockBehaviour func(sub *servicemocks.MockSubscription, a args)

	testCases := []struct {
		testName      string
		args          args
		mockBehaviour mockBehaviour
		inputId       int
		inputBody     string
		expectCode    int
	}{
		{
			testName: "correct test",
			args: args{
				ctx: context.Background(),
				id:  1,
				input: service.SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     ptr(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().Update(a.ctx, a.id, a.input).Return(nil)
			},
			inputId:    1,
			inputBody:  `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "04-2025", "end_date": "06-2025"}`,
			expectCode: http.StatusOK,
		},
		{
			testName: "not found",
			args: args{
				ctx: context.Background(),
				id:  2,
				input: service.SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     ptr(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().Update(a.ctx, a.id, a.input).Return(service.ErrSubscriptionNotFound)
			},
			inputId:    2,
			inputBody:  `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "04-2025", "end_date": "06-2025"}`,
			expectCode: http.StatusNotFound,
		},
		{
			testName: "unexpected error",
			args: args{
				ctx: context.Background(),
				id:  2,
				input: service.SubscriptionInput{
					ServiceName: "Yandex",
					Price:       1000,
					UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
					StartDate:   time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
					EndDate:     ptr(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
			mockBehaviour: func(sub *servicemocks.MockSubscription, a args) {
				sub.EXPECT().Update(a.ctx, a.id, a.input).Return(errors.New("some error"))
			},
			inputId:    2,
			inputBody:  `{"service_name": "Yandex", "price": 1000, "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", "start_date": "04-2025", "end_date": "06-2025"}`,
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			sub := servicemocks.NewMockSubscription(ctrl)
			tc.mockBehaviour(sub, tc.args)

			e := echo.New()
			e.Validator = validator.NewValidator()
			NewRouter(e, &service.Services{Subscription: sub})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/api/v1/subscription/"+strconv.Itoa(tc.inputId), bytes.NewBufferString(tc.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, req)

			assert.Equal(t, tc.expectCode, w.Code)
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
