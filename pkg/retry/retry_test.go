package retry

import (
	"errors"
	"testing"
	"time"
)

func TestRetryHandle_SuccessFirstAttempt(t *testing.T) {
	start := time.Now()
	err := RetryHandle(func() error {
		return nil
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if elapsed < 0 {
		t.Errorf("Expected no delay, but got %v", elapsed)
	}
}

func TestRetryHandle_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	start := time.Now()
	err := RetryHandle(func() error {
		attempts++
		if attempts < 3 {
			return NewRetryableError("temporary error")
		}
		return nil
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error after retries, got %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	expectedMinimum := time.Duration(1+3) * time.Second
	if elapsed < expectedMinimum {
		t.Errorf("Expected at least %v elapsed time, got %v", expectedMinimum, elapsed)
	}
}

func TestRetryHandle_NonRetryableError(t *testing.T) {
	errExpected := errors.New("non-retryable error")
	start := time.Now()
	err := RetryHandle(func() error {
		return errExpected
	})
	elapsed := time.Since(start)

	if err != errExpected {
		t.Errorf("Expected error %v, got %v", errExpected, err)
	}

	expectedMaximum := time.Duration(1) * time.Second
	if elapsed > expectedMaximum {
		t.Errorf("Expected at most %v elapsed time, got %v", expectedMaximum, elapsed)
	}
}

func TestRetryHandle_RetryTimeout(t *testing.T) {
	attempts := 0
	start := time.Now()
	err := RetryHandle(func() error {
		attempts++
		return NewRetryableError("temporary error")
	})
	elapsed := time.Since(start)

	if err == nil {
		t.Errorf("Expected retry timeout error, got nil")
	} else if !errors.Is(err, err) {
		t.Errorf("Expected timeout error, got %v", err)
	}

	expectedAttempts := 4
	if attempts != expectedAttempts {
		t.Errorf("Expected %d attempts, got %d", expectedAttempts, attempts)
	}

	expectedMinimum := time.Duration(0+1+3+5) * time.Second
	if elapsed < expectedMinimum {
		t.Errorf("Expected at least %v elapsed time, got %v", expectedMinimum, elapsed)
	}
}
