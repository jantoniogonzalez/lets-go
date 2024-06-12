package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matchin record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)
