package models

import (
	"errors"
	"fmt"
	"time"
)

//go:generate easyjson -all models_json.go

// CreateSubscriptionJSON Содержит запрос для создания подписки
type CreateSubscriptionJSON struct {
	ServiceName string    `json:"service_name" example:"Yandex Plus" validate:"required"`
	Price       int       `json:"price" example:"400" validate:"required"`
	UserID      string    `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" validate:"required"`
	StartDate   *DateJSON `json:"start_date" example:"07-2025" validate:"required"`
	EndDate     *DateJSON `json:"end_date,omitempty" example:"08-2025"`
}

// UpdateSubscriptionJSON Содержит запрос для обновления подписки
type UpdateSubscriptionJSON struct {
	ID          string    `json:"id" validate:"required"`
	ServiceName string    `json:"service_name,omitempty" example:"Yandex Plus"`
	Price       int       `json:"price,omitempty" example:"400"`
	UserID      string    `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   *DateJSON `json:"start_date,omitempty" example:"07-2025"`
	EndDate     *DateJSON `json:"end_date,omitempty" example:"08-2025"`
}

// Subscriptions Содержит ответ на результат чтения данных
//
//easyjson:json
type Subscriptions []ReadSubscriptionJSON

// ReadSubscriptionJSON Содержит запись результата чтения данных
type ReadSubscriptionJSON struct {
	ServiceName string    `json:"service_name" example:"Yandex Plus" db:"service_name"`
	Price       int       `json:"price" example:"400" db:"sum"`
	UserID      string    `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" db:"user_uuid"`
	StartDate   *DateJSON `json:"start_date" example:"07-2025" db:"begin_date"`
	EndDate     *DateJSON `json:"end_date,omitempty" example:"08-2025" db:"end_date"`
}

// SumSubscriptionsJSON Содержит запрос для подсчета суммарной стоимости всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки
type SumSubscriptionsJSON struct {
	Period      *DateJSON `json:"period" example:"07-2025" validate:"required"`
	ServiceName string    `json:"service_name,omitempty" example:"Yandex Plus"`
	UserID      string    `json:"user_id,omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
}

type DateJSON time.Time

func (d *DateJSON) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("not a json string")
	}

	b = b[1 : len(b)-1]

	t, err := time.Parse("01-2006", string(b))
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	*d = DateJSON(t)

	return nil
}

func (d *DateJSON) MarshalJSON() ([]byte, error) {

	t := time.Time(*d)
	formatted := t.Format("01-2006")

	jsonStr := "\"" + formatted + "\""
	return []byte(jsonStr), nil
}
