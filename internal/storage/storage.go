package storage

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")

	ErrInvalidInput = errors.New("invalid input")
	ErrDatabase     = errors.New("database error")
	ErrConflict     = errors.New("conflict")
)
