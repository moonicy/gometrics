// Package database предоставляет обертку над sql.DB с возможностью повторных попыток при ошибках соединения.
package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/moonicy/gometrics/internal/config"
	"github.com/moonicy/gometrics/pkg/retry"
)

// RetryableDB представляет базу данных с возможностью повторных попыток при ошибках соединения.
type RetryableDB struct {
	db  *sql.DB
	log *zap.SugaredLogger
}

// NewDatabase создает новое соединение с базой данных и возвращает RetryableDB, функцию для закрытия соединения и ошибку, если она произошла.
func NewDatabase(logger *zap.SugaredLogger, cfg config.ServerConfig) (*RetryableDB, func() error, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDsn)
	if err != nil {
		return nil, nil, err
	}
	return &RetryableDB{db: db, log: logger}, db.Close, nil
}

// Ping проверяет соединение с базой данных.
func (db *RetryableDB) Ping() error {
	return db.db.Ping()
}

// ExecContext выполняет SQL-запрос без возвращения строк и поддерживает повторные попытки при ошибках соединения.
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

// QueryContext выполняет SQL-запрос и возвращает несколько строк результата, поддерживая повторные попытки при ошибках соединения.
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

// QueryRowContext выполняет SQL-запрос и возвращает одну строку результата, поддерживая повторные попытки при ошибках соединения.
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

// Begin начинает новую транзакцию и возвращает объект sql.Tx.
func (db *RetryableDB) Begin() (tx *sql.Tx, err error) {
	return db.db.Begin()
}
