package retry

// RetryableError представляет ошибку, которую можно повторить.
type RetryableError struct {
	msg string
}

// NewRetryableError создаёт и возвращает новую RetryableError с заданным сообщением.
func NewRetryableError(msg string) *RetryableError {
	return &RetryableError{msg: msg}
}

// Error возвращает сообщение об ошибке.
func (err RetryableError) Error() string {
	return err.msg
}
