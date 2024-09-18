package storage

import "errors"

// ErrNotFound возвращается, когда запрашиваемые данные не найдены.
var ErrNotFound = errors.New("not found")

// ErrNotValid возвращается, когда данные некорректны или невалидны.
var ErrNotValid = errors.New("not valid")

// ErrConflict возвращается при конфликте данных, например, при нарушении уникальности.
var ErrConflict = errors.New("data conflict")
