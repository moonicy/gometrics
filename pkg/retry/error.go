package retry

type RetryableError struct {
	msg string
}

func NewRetryableError(msg string) *RetryableError {
	return &RetryableError{msg: msg}
}

func (err RetryableError) Error() string {
	return err.msg
}
