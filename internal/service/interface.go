package service

import (
	"context"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"github.com/s-turchinskiy/EffectiveMobile/internal/repository"
	"time"
)

type Servicer interface {
	CreateSubscription(context.Context, models.CreateSubscription) error
	GetSubscriptions(context.Context) (models.Subscriptions, error)
	UpdateSubscription(context.Context, models.UpdateSubscription) error
	DeleteSubscription(ctx context.Context, id uint64) error
	SumSubscriptions(ctx context.Context, data models.SumSubscriptions) (uint64, error)
}

type Service struct {
	Repository    repository.Repository
	retryStrategy []time.Duration
}

func New(rep repository.Repository, retryStrategy []time.Duration) *Service {

	if len(retryStrategy) == 0 {
		retryStrategy = []time.Duration{0}
	}
	return &Service{
		Repository:    rep,
		retryStrategy: retryStrategy,
	}
}
