package storage

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotValid = errors.New("not valid")

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (dbs *DBStorage) Init(ctx context.Context) error {
	_, err := dbs.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS gauge (id serial PRIMARY KEY, name text UNIQUE, value double precision)`)
	if err != nil {
		return err
	}
	_, err = dbs.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS counter (id serial PRIMARY KEY, name text UNIQUE, value int)`)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) SetGauge(ctx context.Context, key string, value float64) error {
	_, err := dbs.db.ExecContext(ctx, `INSERT INTO gauge (name, value) VALUES ($1, $2)
						ON CONFLICT (name) DO UPDATE SET value = $2`, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) AddCounter(ctx context.Context, key string, value int64) error {
	_, err := dbs.db.ExecContext(ctx, `INSERT INTO counter (name, value) VALUES ($1, $2)
						ON CONFLICT (name) DO UPDATE SET value = counter.value + EXCLUDED.value`, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	row := dbs.db.QueryRowContext(ctx, `SELECT value FROM counter WHERE name = $1`, key)
	var value sql.NullInt64

	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	if value.Valid {
		return value.Int64, nil
	}
	return 0, ErrNotValid
}

func (dbs *DBStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	row := dbs.db.QueryRowContext(ctx, `SELECT value FROM gauge WHERE name = $1`, key)
	var value sql.NullFloat64

	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	if value.Valid {
		return value.Float64, nil
	}
	return 0, ErrNotValid
}

func (dbs *DBStorage) GetMetrics(ctx context.Context) (map[string]int64, map[string]float64, error) {
	rowsGauge, err := dbs.db.QueryContext(ctx, `SELECT name, value FROM gauge ORDER BY name`)
	if err != nil {
		return nil, nil, err
	}
	defer rowsGauge.Close()

	gauge := make(map[string]float64)

	for rowsGauge.Next() {
		var name string
		var value float64
		err = rowsGauge.Scan(&name, &value)
		if err != nil {
			return nil, nil, err
		}

		gauge[name] = value
	}

	err = rowsGauge.Err()
	if err != nil {
		return nil, nil, err
	}

	rowsCounter, err := dbs.db.Query(`SELECT name, value FROM counter ORDER BY name`)
	if err != nil {
		return nil, nil, err
	}
	defer rowsCounter.Close()

	counter := make(map[string]int64)

	for rowsCounter.Next() {
		var name string
		var value int64
		err = rowsCounter.Scan(&name, &value)
		if err != nil {
			return nil, nil, err
		}

		counter[name] = value
	}

	err = rowsCounter.Err()
	if err != nil {
		return nil, nil, err
	}

	return counter, gauge, nil
}
