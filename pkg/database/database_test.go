// database_test.go
package database

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap/zaptest"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupMockDB(t *testing.T) (*RetryableDB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании sqlmock: %v", err)
	}
	logger := zaptest.NewLogger(t).Sugar()
	retryableDB := &RetryableDB{
		db:  db,
		log: logger,
	}
	closeFunc := func() {
		_ = db.Close()
	}
	return retryableDB, mock, closeFunc
}

func TestRetryableDB_Ping_Success(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	mock.ExpectPing()

	err := retryableDB.Ping()
	if err != nil {
		t.Errorf("Ожидали успешный Ping, получили ошибку: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_ExecContext_Success(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	query := "INSERT INTO metrics \\(name, value\\) VALUES \\(\\$1, \\$2\\)"
	mock.ExpectExec(query).
		WithArgs("cpu", 0.75).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	result, err := retryableDB.ExecContext(ctx, "INSERT INTO metrics (name, value) VALUES ($1, $2)", "cpu", 0.75)
	if err != nil {
		t.Errorf("Ожидали успешный ExecContext, получили ошибку: %v", err)
	}
	if result == nil {
		t.Errorf("Ожидали валидный sql.Result, получили nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_ExecContext_NonRetryableError(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	query := "INSERT INTO metrics \\(name, value\\) VALUES \\(\\$1, \\$2\\)"
	nonRetryableErr := errors.New("syntax error")
	mock.ExpectExec(query).
		WithArgs("cpu", 0.75).
		WillReturnError(nonRetryableErr)

	ctx := context.Background()
	result, err := retryableDB.ExecContext(ctx, "INSERT INTO metrics (name, value) VALUES ($1, $2)", "cpu", 0.75)
	if err == nil {
		t.Errorf("Ожидали ошибку при ExecContext, но ошибки нет")
	}
	if !errors.Is(err, nonRetryableErr) {
		t.Errorf("Ожидали ошибку '%v', получили '%v'", nonRetryableErr, err)
	}
	if result != nil {
		t.Errorf("Ожидали sql.Result равным nil, получили %v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_QueryContext_Success(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	query := "SELECT name, value FROM metrics WHERE name = \\$1"
	rows := sqlmock.NewRows([]string{"name", "value"}).
		AddRow("cpu", 0.75)
	mock.ExpectQuery(query).
		WithArgs("cpu").
		WillReturnRows(rows)

	ctx := context.Background()
	resultRows, err := retryableDB.QueryContext(ctx, "SELECT name, value FROM metrics WHERE name = $1", "cpu")
	if err != nil {
		t.Errorf("Ожидали успешный QueryContext, получили ошибку: %v", err)
	}
	if resultRows == nil {
		t.Errorf("Ожидали *sql.Rows, получили nil")
	}
	if resultRows.Err() != nil {
		t.Errorf("Ожидали успешный QueryContext, получили ошибку в Rows: %v", resultRows.Err())
	}
	defer resultRows.Close()

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_QueryRowContext_Success(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	query := "SELECT name, value FROM metrics WHERE name = \\$1"
	rows := sqlmock.NewRows([]string{"name", "value"}).
		AddRow("cpu", 0.75)
	mock.ExpectQuery(query).
		WithArgs("cpu").
		WillReturnRows(rows)

	ctx := context.Background()
	row := retryableDB.QueryRowContext(ctx, "SELECT name, value FROM metrics WHERE name = $1", "cpu")

	var name string
	var value float64
	err := row.Scan(&name, &value)
	if err != nil {
		t.Errorf("Ожидали успешный Scan, получили ошибку: %v", err)
	}
	if name != "cpu" || value != 0.75 {
		t.Errorf("Ожидали данные ('cpu', 0.75), получили ('%s', %v)", name, value)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_QueryRowContext_NonRetryableError(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	query := "SELECT name, value FROM metrics WHERE name = \\$1"
	nonRetryableErr := errors.New("syntax error")
	mock.ExpectQuery(query).
		WithArgs("cpu").
		WillReturnError(nonRetryableErr)

	ctx := context.Background()
	row := retryableDB.QueryRowContext(ctx, "SELECT name, value FROM metrics WHERE name = $1", "cpu")

	var name string
	var value float64
	err := row.Scan(&name, &value)
	if err == nil {
		t.Errorf("Ожидали ошибку при Scan, но ошибки нет")
	}
	if !errors.Is(err, nonRetryableErr) {
		t.Errorf("Ожидали ошибку '%v', получили '%v'", nonRetryableErr, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_Begin_Success(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	mock.ExpectBegin()

	tx, err := retryableDB.Begin()
	if err != nil {
		t.Errorf("Ожидали успешный Begin, получили ошибку: %v", err)
	}
	if tx == nil {
		t.Errorf("Ожидали валидный *sql.Tx, получили nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}

func TestRetryableDB_Begin_Error(t *testing.T) {
	retryableDB, mock, closeFunc := setupMockDB(t)
	defer closeFunc()

	mock.ExpectBegin().WillReturnError(errors.New("begin failed"))

	tx, err := retryableDB.Begin()
	if err == nil {
		t.Errorf("Ожидали ошибку при Begin, но ошибки нет")
	}
	if tx != nil {
		t.Errorf("Ожидали *sql.Tx равным nil, получили %v", tx)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Ожидания не были выполнены: %v", err)
	}
}
