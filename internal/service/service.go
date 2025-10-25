package service

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"time"
)

var (
	ErrEndDateLessStartDate = fmt.Errorf("end date less start date")
	ErrNotRowWithThisID     = fmt.Errorf("not row with this ID")
)

func (s Service) CreateSubscription(ctx context.Context, data models.CreateSubscription) error {

	if data.EndDate != nil && time.Time.Before(*data.EndDate, *data.StartDate) {
		return ErrEndDateLessStartDate
	}
	return s.Repository.CreateSubscription(ctx, data)
}

func (s Service) GetSubscriptions(ctx context.Context) (models.Subscriptions, error) {

	return s.Repository.GetSubscriptions(ctx)
}

func (s Service) DeleteSubscription(ctx context.Context, id uint64) error {

	return s.Repository.DeleteSubscription(ctx, id)
}

func (s Service) UpdateSubscription(ctx context.Context, data models.UpdateSubscription) error {

	if data.EndDate != nil && time.Time.Before(*data.EndDate, *data.StartDate) {
		return ErrEndDateLessStartDate
	}
	return s.Repository.UpdateSubscription(ctx, data)
}

func (s Service) SumSubscriptions(ctx context.Context, data models.SumSubscriptions) (uint64, error) {
	return s.Repository.SumSubscriptions(ctx, data)
}
