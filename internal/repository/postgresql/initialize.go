package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/s-turchinskiy/EffectiveMobile/internal/common/errors"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/repository"
	"go.uber.org/zap"
	"time"
)

type PostgreSQL struct {
	db   *sqlx.DB
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, addr, dbname string) (repository.Repository, error) {

	schemaName := "effectivemobile"
	logger.Log.Debug("addr for Sql.Open: ", addr)

	db, err := sqlx.Open("pgx", addr)
	if err != nil {
		return nil, commonerrors.WrapError(err)
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, commonerrors.WrapError(err)
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(30 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pool, err := pgxpool.New(ctx, addr)
	if err != nil {
		return nil, commonerrors.WrapError(err)
	}

	p := &PostgreSQL{db: db, pool: pool}

	_, err = p.db.ExecContext(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schemaName))
	if err != nil {
		return nil, commonerrors.WrapError(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{SchemaName: schemaName})
	if err != nil {
		return nil, commonerrors.WrapError(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/repository/postgresql/migrations", dbname, driver)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, commonerrors.WrapError(err)
	}

	return p, nil

}

func (p *PostgreSQL) Close(ctx context.Context) error {

	p.pool.Close()

	err := p.db.Close()

	if err != nil {
		logger.Log.Infow("PostgreSQL stopped with error", zap.String("error", err.Error()))
	} else {
		logger.Log.Infow("PostgreSQL stopped")
	}
	return err

}
