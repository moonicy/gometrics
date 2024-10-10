package retry

import (
	"testing"
)

func TestNewRetryableError(t *testing.T) {
	msg := "temporary error"
	err := NewRetryableError(msg)

	if err == nil {
		t.Fatal("Expected RetryableError, got nil")
	}

	if err.msg != msg {
		t.Errorf("Expected message '%s', got '%s'", msg, err.msg)
	}
}

func TestRetryableError_Error(t *testing.T) {
	msg := "another temporary error"
	err := RetryableError{msg: msg}

	if err.Error() != msg {
		t.Errorf("Expected Error() to return '%s', got '%s'", msg, err.Error())
	}
}
