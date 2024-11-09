package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stretchr/testify/assert"
)

// Тест Init: проверка инициализации таблиц
func TestDBStorage_Init(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS gauge").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS counter").WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewDBStorage(db)
	err = storage.Init(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест SetGauge: вставка и обновление метрики типа gauge
func TestDBStorage_SetGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO gauge").WithArgs("gauge1", 10.5).WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewDBStorage(db)
	err = storage.SetGauge(context.Background(), "gauge1", 10.5)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест SetGauge: проверка обработки ошибки ограничения целостности
func TestDBStorage_SetGauge_ConflictError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	mock.ExpectExec("INSERT INTO gauge").WillReturnError(pgErr)

	storage := NewDBStorage(db)
	err = storage.SetGauge(context.Background(), "gauge1", 10.5)
	assert.Equal(t, pgErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест AddCounter: вставка и обновление метрики типа counter
func TestDBStorage_AddCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO counter").WithArgs("counter1", 10).WillReturnResult(sqlmock.NewResult(1, 1))

	storage := NewDBStorage(db)
	err = storage.AddCounter(context.Background(), "counter1", 10)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест AddCounter: проверка обработки ошибки ограничения целостности
func TestDBStorage_AddCounter_ConflictError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pgErr := &pgconn.PgError{Code: pgerrcode.UniqueViolation}
	mock.ExpectExec("INSERT INTO counter").WillReturnError(pgErr)

	storage := NewDBStorage(db)
	err = storage.AddCounter(context.Background(), "counter1", 10)
	assert.Equal(t, pgErr, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест GetCounter: получение значения метрики типа counter
func TestDBStorage_GetCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"}).AddRow(100)
	mock.ExpectQuery("SELECT value FROM counter").WithArgs("counter1").WillReturnRows(rows)

	storage := NewDBStorage(db)
	value, err := storage.GetCounter(context.Background(), "counter1")
	assert.NoError(t, err)
	assert.Equal(t, int64(100), value)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест GetCounter: обработка отсутствующих данных
func TestDBStorage_GetCounter_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT value FROM counter").WithArgs("counter1").WillReturnError(sql.ErrNoRows)

	storage := NewDBStorage(db)
	_, err = storage.GetCounter(context.Background(), "counter1")
	assert.Equal(t, ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест GetGauge: получение значения метрики типа gauge
func TestDBStorage_GetGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"value"}).AddRow(15.5)
	mock.ExpectQuery("SELECT value FROM gauge").WithArgs("gauge1").WillReturnRows(rows)

	storage := NewDBStorage(db)
	value, err := storage.GetGauge(context.Background(), "gauge1")
	assert.NoError(t, err)
	assert.Equal(t, 15.5, value)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест GetGauge: обработка отсутствующих данных
func TestDBStorage_GetGauge_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT value FROM gauge").WithArgs("gauge1").WillReturnError(sql.ErrNoRows)

	storage := NewDBStorage(db)
	_, err = storage.GetGauge(context.Background(), "gauge1")
	assert.Equal(t, ErrNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Тест GetMetrics: получение всех метрик типа counter и gauge
func TestDBStorage_GetMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rowsGauge := sqlmock.NewRows([]string{"name", "value"}).AddRow("gauge1", 10.5).AddRow("gauge2", 20.5)
	rowsCounter := sqlmock.NewRows([]string{"name", "value"}).AddRow("counter1", 100).AddRow("counter2", 200)

	mock.ExpectQuery("SELECT name, value FROM gauge ORDER BY name").WillReturnRows(rowsGauge)
	mock.ExpectQuery("SELECT name, value FROM counter ORDER BY name").WillReturnRows(rowsCounter)

	storage := NewDBStorage(db)
	counters, gauges, err := storage.GetMetrics(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"gauge1": 10.5, "gauge2": 20.5}, gauges)
	assert.Equal(t, map[string]int64{"counter1": 100, "counter2": 200}, counters)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBStorage_SetMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Ожидаемые данные для таблицы counter
	counter := map[string]int64{
		"counter1": 100,
		"counter2": 200,
	}

	// Ожидаемые данные для таблицы gauge
	gauge := map[string]float64{
		"gauge1": 10.5,
		"gauge2": 20.5,
	}

	// Начало транзакции
	mock.ExpectBegin()

	// Проверка вставки или обновления значений для counter
	counterQuery := "INSERT INTO counter"
	mock.ExpectPrepare(counterQuery)
	mock.ExpectExec(counterQuery).
		WithArgs("counter1", 100, "counter2", 200).
		WillReturnResult(sqlmock.NewResult(1, 2))

	// Проверка вставки или обновления значений для gauge
	gaugeQuery := "INSERT INTO gauge"
	mock.ExpectPrepare(gaugeQuery)
	mock.ExpectExec(gaugeQuery).
		WithArgs("gauge1", 10.5, "gauge2", 20.5).
		WillReturnResult(sqlmock.NewResult(1, 2))

	// Завершение транзакции
	mock.ExpectCommit()

	storage := NewDBStorage(db)
	err = storage.SetMetrics(context.Background(), counter, gauge)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
