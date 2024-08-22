package domain

import "github.com/pkg/errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidParams  = errors.New("invalid parameters")
	ErrInvalidCost    = errors.New("cost must be a multiple of five")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidToken      = errors.New("unauthorized user")
	ErrUnauthorized      = errors.New("unauthorized user")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrPermissionDenied  = errors.New("permission denied")

	ErrInvalidCoin                = errors.New("invalid coin, use 5, 10, 20, 50, 100 cent coins")
	ErrProductNotFound            = errors.New("product not found")
	ErrInsufficientBalance        = errors.New("insufficient balance")
	ErrInsufficientProductsAmount = errors.New("insufficient products amount")
)
