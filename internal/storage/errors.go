package storage

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotValid = errors.New("not valid")
var ErrConflict = errors.New("data conflict")
