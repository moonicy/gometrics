package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/retry"
	"go.uber.org/zap"
)

type RetryableDB struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

func NewDatabase(logger *zap.SugaredLogger, cfg config.ServerConfig) (*RetryableDB, func() error, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		return nil, nil, err
	}
	return &RetryableDB{db: db, log: logger}, db.Close, nil
}

func (db *RetryableDB) Ping() error {
	return db.db.Ping()
}

func (db *RetryableDB) ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	db.log.Info("opening database")
	err = retry.RetryHandle(func() error {
		result, err = db.db.ExecContext(ctx, query, args...)
		if err == nil {
			return nil
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsConnectionException(pgErr.Code) {
				db.log.Warnw("database connection error", "error", err)
				return retry.NewRetryableError(pgErr.Message)
			}
		}
		var connErr *pgconn.ConnectError
		if errors.As(err, &connErr) {
			db.log.Warnw("database connection error", "error", err)
			return retry.NewRetryableError(err.Error())
		}
		db.log.Errorw("database retry error", "error", err)
		return err
	})
	if err != nil {
		db.log.Errorw("database error", "error", err)
		return nil, err
	}
	db.log.Info("database opened")
	return result, nil
}

func (db *RetryableDB) QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	db.log.Info("opening database")
	err = retry.RetryHandle(func() error {
		rows, err = db.db.QueryContext(ctx, query, args...)
		if err == nil {
			return nil
		}
		if rows.Err() != nil {
			err = rows.Err()
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			db.log.Warnw("database connection error", "error", err)
			return retry.NewRetryableError(pgErr.Message)
		}
		var connErr *pgconn.ConnectError
		if errors.As(err, &connErr) {
			db.log.Warnw("database connection error", "error", err)
			return retry.NewRetryableError(err.Error())
		}
		db.log.Errorw("database retry error", "error", err)
		return err
	})
	if err != nil {
		db.log.Errorw("database error", "error", err)
		return nil, err
	}
	db.log.Info("database opened")
	return rows, nil
}

func (db *RetryableDB) QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row) {
	var err error
	db.log.Info("opening database")
	err = retry.RetryHandle(func() error {
		row = db.db.QueryRowContext(ctx, query, args...)
		err = row.Err()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
			db.log.Warnw("database connection error", "error", err)
			return retry.NewRetryableError(pgErr.Message)
		}
		var connErr *pgconn.ConnectError
		if errors.As(err, &connErr) {
			db.log.Warnw("database connection error", "error", err)
			return retry.NewRetryableError(err.Error())
		}
		db.log.Errorw("database retry error", "error", err)
		return err
	})
	if err != nil {
		db.log.Errorw("database error", "error", err)
		return row
	}
	db.log.Info("database opened")
	return row
}

func (db *RetryableDB) Begin() (tx *sql.Tx, err error) {
	return db.db.Begin()
}
