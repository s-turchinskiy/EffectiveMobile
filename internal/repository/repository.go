package repository

import (
	"context"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
)

type Repository interface {
	Close(ctx context.Context) error
	CreateSubscription(context.Context, models.CreateSubscription) error
	GetSubscriptions(context.Context) ([]models.ReadSubscriptionJSON, error)
	UpdateSubscription(context.Context, models.UpdateSubscription) error
	DeleteSubscription(ctx context.Context, id uint64) error
	SumSubscriptions(ctx context.Context, data models.SumSubscriptions) (uint64, error)
}
