package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/s-turchinskiy/EffectiveMobile/internal/common/common"
	commonerrors "github.com/s-turchinskiy/EffectiveMobile/internal/common/errors"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"github.com/s-turchinskiy/EffectiveMobile/internal/service"
	"strings"
	"time"
)

func (p *PostgreSQL) CreateSubscription(ctx context.Context, data models.CreateSubscription) error {

	request1, err := getRequest("create_user.sql")
	if err != nil {
		return err
	}

	request2, err := getRequest("create_service.sql")
	if err != nil {
		return err
	}

	request3, err := getRequest("create_subscription.sql")
	if err != nil {
		return err
	}

	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquire connection: %w", err)
	}
	defer conn.Release()

	batch := new(pgx.Batch)

	batch.Queue(request1, data.UserID)
	batch.Queue(request2, data.ServiceName)
	var endDate time.Time
	if data.EndDate != nil {
		endDate = *data.EndDate
	}
	batch.Queue(request3, data.UserID, data.ServiceName, data.Price, *data.StartDate, common.Ternary(data.EndDate == nil, nil, endDate))

	if batch.Len() == 0 {
		return nil
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return commonerrors.WrapError(err)
	}

	result := tx.SendBatch(ctx, batch)

	_, err = result.Exec()
	if err != nil {
		return commonerrors.WrapError(err)
	}

	err = result.Close()
	if err != nil {
		_ = tx.Rollback(ctx)

		if commonerrors.IsDuplicateKeyError(err) {
			return commonerrors.ErrDuplicateKey
		}

		return commonerrors.WrapError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return commonerrors.WrapError(err)
	}

	return nil

}

func (p *PostgreSQL) GetSubscriptions(ctx context.Context) ([]models.ReadSubscriptionJSON, error) {

	request, err := getRequest("get_subscriptions.sql")
	if err != nil {
		return nil, err
	}

	var result []models.ReadSubscriptionJSON
	err = p.db.SelectContext(ctx, &result, request)

	if err != nil {
		return nil, commonerrors.WrapError(err)
	}

	return result, nil

}

func (p *PostgreSQL) DeleteSubscription(ctx context.Context, id uint64) error {

	requestSelect := "SELECT id FROM effectivemobile.subscriptions WHERE id = $1"
	requestDelete := "DELETE FROM effectivemobile.subscriptions WHERE id = $1"

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, requestSelect, id)
	err = row.Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return service.ErrNotRowWithThisID
	case err != nil:
		return commonerrors.WrapError(err)
	}

	_, err = tx.Exec(requestDelete, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (p *PostgreSQL) UpdateSubscription(ctx context.Context, subscription models.UpdateSubscription) error {

	requestSelect := "SELECT id FROM effectivemobile.subscriptions WHERE id = $1"

	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, requestSelect, subscription.ID)
	err = row.Scan(&subscription.ID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return service.ErrNotRowWithThisID
	case err != nil:
		return commonerrors.WrapError(err)
	}

	requestUpdate, err := getRequestUpdate(subscription)
	if err != nil {
		return err
	}

	_, err = tx.Exec(requestUpdate, subscription.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (p *PostgreSQL) SumSubscriptions(ctx context.Context, data models.SumSubscriptions) (uint64, error) {

	request := getRequestSumSubscriptions(data)

	var result uint64
	row := p.db.QueryRowContext(ctx, request, data.Period)
	err := row.Scan(&result)

	if err != nil {
		return 0, commonerrors.WrapError(err)
	}

	return result, nil

}

func getRequestUpdate(subscription models.UpdateSubscription) (string, error) {

	requestUpdate := "UPDATE effectivemobile.subscriptions SET $2 WHERE id = $1"

	var values []string
	if subscription.UserID != "" {
		s := fmt.Sprintf("user_id = (select id from effectivemobile.users where uuid = '%s')", subscription.UserID)
		values = append(values, s)
	}

	if subscription.ServiceName != "" {
		s := fmt.Sprintf("service_id = (select id from effectivemobile.services where name = '%s')", subscription.ServiceName)
		values = append(values, s)
	}

	if subscription.Price != 0 {
		s := fmt.Sprintf("sum = '%d'", subscription.Price)
		values = append(values, s)
	}

	if subscription.StartDate != nil {

		s := fmt.Sprintf("begin_date = '%s'", subscription.StartDate.Format(time.DateTime))
		values = append(values, s)
	}

	if subscription.EndDate != nil {
		s := fmt.Sprintf("end_date = '%s'", subscription.EndDate.Format(time.DateTime))
		values = append(values, s)
	}

	if len(values) == 0 {
		return "", fmt.Errorf("no data for update")
	}

	return strings.Replace(requestUpdate, "$2", strings.Join(values, ","), 1), nil

}

func getRequestSumSubscriptions(subscription models.SumSubscriptions) string {

	request := `select coalesce(
	(select
	sum(s.sum)
	from
	effectivemobile.subscriptions s
	where
	s.begin_date <= $1
	and (s.end_date is null
	or s.end_date >= $1) $2), 0) as sum`

	var requestWhere string
	if subscription.UserID != "" {
		requestWhere = requestWhere + fmt.Sprintf("and s.user_id = (select id from effectivemobile.users where uuid = '%s')", subscription.UserID)
	}

	if subscription.ServiceName != "" {
		requestWhere = requestWhere + fmt.Sprintf("and s.service_id = (select id from effectivemobile.services where name = '%s')", subscription.ServiceName)
	}

	return strings.Replace(request, "$2", requestWhere, 1)

}
