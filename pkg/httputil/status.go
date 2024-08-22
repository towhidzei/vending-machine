package httputil

import (
	"net/http"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/labstack/echo"
)

func StatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if httpErr, ok := err.(*echo.HTTPError); ok {
		return httpErr.Code
	}
	switch err {
	case domain.ErrWrongCredentials, domain.ErrInvalidToken, domain.ErrUnauthorized:
		return http.StatusUnauthorized
	case domain.ErrPermissionDenied:
		return http.StatusForbidden
	case domain.ErrInvalidParams:
		return http.StatusBadRequest
	case domain.ErrUserNotFound, domain.ErrProductNotFound:
		return http.StatusNotFound
	case domain.ErrInsufficientProductsAmount, domain.ErrInsufficientBalance:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
