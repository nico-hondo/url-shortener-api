package apperror

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrCredentialsEmpty   = errors.New("email or password cannot be empty")

	// Generic data access related
	ErrNotFound       = errors.New("record not found")
	ErrNoRowsAffected = errors.New("no rows affected")

	// General
	ErrInternal     = errors.New("internal server error")
	ErrUnauthorized = errors.New("unauthorized")

	ErrMissingKey = errors.New("jwt key not found")

	// pkg/hash.go

	ErrInvalidHash               = errors.New("argon2id: hash is not in the correct format")
	ErrIncompatibleVariant       = errors.New("argon2id: incompatible variant of argon2")
	ErrIncompatibleVersion       = errors.New("argon2id: incompatible version of argon2")
	ErrMismatchedHashAndPassword = errors.New("argon2id: the provided password does not match the hash")
)
