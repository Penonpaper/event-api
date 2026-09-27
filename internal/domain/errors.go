package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")

	ErrDataBase = errors.New("database error")

	//ErrFieldIsEmpty = errors.New("field is empty")

	// 401 - Unauthorized
	ErrRefreshToken = errors.New("invalid refresh token")
	ErrAccessToken  = errors.New("invalid access token")

	ErrDuplicateEmail = errors.New("duplicate email")
	ErrInternal       = errors.New("internal error")

	ErrRedis = errors.New("redis error")

	ErrConditions = errors.New("failed condition")
)
