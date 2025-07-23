package pgdb

import (
	"subscription_service/internal/model/dbmodel"
	"subscription_service/internal/repo/pgerrs"
	"testing"
	"time"
)

func (s *pgdbTestSuite) TestSubscriptionRepo_Create() {
	testCases := []struct {
		testName  string
		sub       dbmodel.Subscription
		expectErr error
	}{
		{
			testName: "correct test",
			sub: dbmodel.Subscription{
				ServiceName: "Yandex",
				Price:       1000,
				UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
				StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     ptr(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
			expectErr: nil,
		},
		{
			testName: "correct test with null end date",
			sub: dbmodel.Subscription{
				ServiceName: "Google",
				Price:       500,
				UserId:      "2234696a-d069-4fad-a3ed-f27c13651c3a",
				StartDate:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:     nil,
			},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			err := s.sub.Create(s.ctx, tc.sub)

			s.Assert().Equal(tc.expectErr, err)

			if tc.expectErr == nil {
				sql, args, _ := s.pg.Builder.
					Select("service_name", "price", "user_id", "start_date", "end_date").
					From(subscriptionTable).
					Where("user_id = ?", tc.sub.UserId).
					ToSql()

				var actual dbmodel.Subscription

				err = s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan(
					&actual.ServiceName,
					&actual.Price,
					&actual.UserId,
					&actual.StartDate,
					&actual.EndDate,
				)
				s.Assert().NoError(err)
				s.Assert().Equal(tc.sub, actual)
			}
		})
	}
}

func (s *pgdbTestSuite) TestSubscriptionRepo_FindById() {
	defaultSub := dbmodel.Subscription{
		ServiceName: "Yandex",
		Price:       1000,
		UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
		StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     ptr(time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)),
	}

	sql, args, _ := s.pg.Builder.
		Insert(subscriptionTable).
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(defaultSub.ServiceName, defaultSub.Price, defaultSub.UserId, defaultSub.StartDate, defaultSub.EndDate).
		Suffix("RETURNING id").
		ToSql()

	var defaultId int

	if err := s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan(&defaultId); err != nil {
		panic(err)
	}

	testCases := []struct {
		testName     string
		id           int
		expectOutput dbmodel.Subscription
		expectErr    error
	}{
		{
			testName:     "correct test",
			id:           defaultId,
			expectOutput: defaultSub,
			expectErr:    nil,
		},
		{
			testName:  "not found",
			id:        0,
			expectErr: pgerrs.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			sub, err := s.sub.FindById(s.ctx, tc.id)

			s.Assert().Equal(tc.expectErr, err)
			s.Assert().Equal(tc.expectOutput, sub)
		})
	}
}

func (s *pgdbTestSuite) TestSubscriptionRepo_FindPrice() {
	subscriptions := []dbmodel.Subscription{
		{
			ServiceName: "Yandex",
			Price:       500,
			UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
			StartDate:   time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     ptr(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			ServiceName: "Yandex",
			Price:       500,
			UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
			StartDate:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     nil,
		},
		{
			ServiceName: "Google",
			Price:       1000,
			UserId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
			StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     ptr(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			ServiceName: "Google",
			Price:       1000,
			UserId:      "2344696a-d069-4fad-a3ed-f27c13651c3a",
			StartDate:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     ptr(time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			ServiceName: "VK",
			Price:       400,
			UserId:      "2344696a-d069-4fad-a3ed-f27c13651c3a",
			StartDate:   time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
			EndDate:     ptr(time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	for _, sub := range subscriptions {
		sql, args, _ := s.pg.Builder.
			Insert(subscriptionTable).
			Columns("service_name", "price", "user_id", "start_date", "end_date").
			Values(sub.ServiceName, sub.Price, sub.UserId, sub.StartDate, sub.EndDate).
			ToSql()

		if _, err := s.pg.Pool.Exec(s.ctx, sql, args...); err != nil {
			panic(err)
		}
	}

	testCases := []struct {
		testName    string
		service     string
		userId      string
		start       time.Time
		end         time.Time
		expectPrice int
	}{
		{
			testName:    "find by time interval (all)",
			service:     "",
			userId:      "",
			start:       time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 3400,
		},
		{
			testName:    "find by time interval (3,4)",
			service:     "",
			userId:      "",
			start:       time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 2500, // Yandex (потому что без срока окончания) + 2 * Google
		},
		{
			testName:    "find by time interval (2)",
			service:     "",
			userId:      "",
			start:       time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 500, // Yandex (2) потому что без срока
		},
		{
			testName:    "find by service (Yandex)",
			service:     "Yandex",
			userId:      "",
			start:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 1000,
		},
		{
			testName:    "find by service (Google)",
			service:     "Google",
			userId:      "",
			start:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 2000,
		},
		{
			testName:    "find by user",
			service:     "",
			userId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
			start:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 2000,
		},
		{
			testName:    "find by user and service",
			service:     "Yandex",
			userId:      "6114696a-d069-4fad-a3ed-f27c13651c3a",
			start:       time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 1000,
		},
		{
			testName:    "empty search range",
			service:     "",
			userId:      "",
			start:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			end:         time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			expectPrice: 0,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			price, err := s.sub.FindPrice(s.ctx, tc.service, tc.userId, tc.start, tc.end)

			s.Assert().NoError(err)

			s.Assert().Equal(tc.expectPrice, price)
		})
	}
}

func (s *pgdbTestSuite) TestSubscriptionRepo_Delete() {
	sql, args, _ := s.pg.Builder.
		Insert(subscriptionTable).
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values("Yandex", 100, "6114696a-d069-4fad-a3ed-f27c13651c3a", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil).
		Suffix("RETURNING id").
		ToSql()

	var defaultId int

	if err := s.pg.Pool.QueryRow(s.ctx, sql, args...).Scan(&defaultId); err != nil {
		panic(err)
	}

	testCases := []struct {
		testName  string
		id        int
		expectErr error
	}{
		{
			testName:  "subscription does not exists",
			id:        0,
			expectErr: pgerrs.ErrNotFound,
		},
		{
			testName:  "correct test",
			id:        defaultId,
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.testName, func(t *testing.T) {
			err := s.sub.Delete(s.ctx, tc.id)

			s.Assert().Equal(tc.expectErr, err)

			if tc.expectErr == nil {
				sql = "SELECT EXISTS(SELECT id FROM subscription WHERE id = $1)"

				var status bool

				err = s.pg.Pool.QueryRow(s.ctx, sql, tc.id).Scan(&status)
				s.Assert().NoError(err)

				s.Assert().False(status)
			}
		})
	}
}

func ptr[T any](t T) *T {
	return &t
}
