package models

import "time"

type CreateSubscription struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   *time.Time
	EndDate     *time.Time
}

type UpdateSubscription struct {
	ID          uint64
	ServiceName string
	Price       int
	UserID      string
	StartDate   *time.Time
	EndDate     *time.Time
}

type SumSubscriptions struct {
	Period      *time.Time
	ServiceName string
	UserID      string
}
