package dbmodel

import "time"

type Subscription struct {
	Id          int
	ServiceName string
	Price       int
	UserId      string
	StartDate   time.Time
	EndDate     *time.Time
}
