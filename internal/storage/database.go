package storage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// DB определяет интерфейс для взаимодействия с базой данных.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
	Begin() (tx *sql.Tx, err error)
}

// DBStorage представляет хранилище метрик, использующее базу данных.
type DBStorage struct {
	db DB
}

// NewDBStorage создаёт и возвращает новое хранилище метрик с использованием заданного интерфейса DB.
func NewDBStorage(db DB) *DBStorage {
	return &DBStorage{db: db}
}

// Init инициализирует хранилище, создавая необходимые таблицы в базе данных.
func (dbs *DBStorage) Init(ctx context.Context) error {
	_, err := dbs.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS gauge (id serial PRIMARY KEY, name text UNIQUE, value double precision)`)
	if err != nil {
		return err
	}
	_, err = dbs.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS counter (id serial PRIMARY KEY, name text UNIQUE, value bigint)`)
	if err != nil {
		return err
	}
	return nil
}

// SetGauge устанавливает значение метрики типа gauge с заданным именем и значением.
func (dbs *DBStorage) SetGauge(ctx context.Context, key string, value float64) error {
	_, err := dbs.db.ExecContext(ctx, `INSERT INTO gauge (name, value) VALUES ($1, $2)
						ON CONFLICT (name) DO UPDATE SET value = $2`, key, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}
	return nil
}

// AddCounter увеличивает значение метрики типа counter с заданным именем на указанное значение.
func (dbs *DBStorage) AddCounter(ctx context.Context, key string, value int64) error {
	_, err := dbs.db.ExecContext(ctx, `INSERT INTO counter (name, value) VALUES ($1, $2)
						ON CONFLICT (name) DO UPDATE SET value = counter.value + EXCLUDED.value`, key, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = ErrConflict
		}
		return err
	}
	return nil
}

// GetCounter возвращает текущее значение метрики типа counter с заданным именем.
func (dbs *DBStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	row := dbs.db.QueryRowContext(ctx, `SELECT value FROM counter WHERE name = $1`, key)
	var value sql.NullInt64

	err := row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	if value.Valid {
		return value.Int64, nil
	}
	return 0, ErrNotValid
}

// GetGauge возвращает текущее значение метрики типа gauge с заданным именем.
func (dbs *DBStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	row := dbs.db.QueryRowContext(ctx, `SELECT value FROM gauge WHERE name = $1`, key)
	var value sql.NullFloat64

	err := row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	if value.Valid {
		return value.Float64, nil
	}
	return 0, ErrNotValid
}

// GetMetrics возвращает все сохранённые метрики типа counter и gauge.
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

	rowsCounter, err := dbs.db.QueryContext(ctx, `SELECT name, value FROM counter ORDER BY name`)
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

// SetMetrics сохраняет переданные метрики типа counter и gauge в базе данных.
func (dbs *DBStorage) SetMetrics(ctx context.Context, counter map[string]int64, gauge map[string]float64) error {
	tx, err := dbs.db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dbs.setCounters(ctx, tx, counter)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = dbs.setGauges(ctx, tx, gauge)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dbs *DBStorage) setCounters(ctx context.Context, tx *sql.Tx, counter map[string]int64) error {
	sqlStr := "INSERT INTO counter(name, value) VALUES "
	vals := make([]interface{}, 0, len(counter))

	n := 0
	for name, value := range counter {
		sqlStr += "($" + strconv.Itoa(n+1) + ", $" + strconv.Itoa(n+2) + "),"
		n += 2
		vals = append(vals, name, value)
	}
	// trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += "ON CONFLICT (name) DO UPDATE SET value = counter.value + EXCLUDED.value"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) setGauges(ctx context.Context, tx *sql.Tx, gauge map[string]float64) error {
	sqlStr := "INSERT INTO gauge(name, value) VALUES "
	vals := make([]interface{}, 0, len(gauge))

	n := 0
	for name, value := range gauge {
		sqlStr += "($" + strconv.Itoa(n+1) + ", $" + strconv.Itoa(n+2) + "),"
		n += 2
		vals = append(vals, name, value)
	}
	// trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += "ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value"

	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		return err
	}
	return nil
}
