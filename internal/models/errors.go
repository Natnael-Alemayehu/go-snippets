package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Error used if a user tries to login with incorrect username or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// Error for duplicate email
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
