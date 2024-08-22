package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/labstack/echo"
)

func (m *UserMiddleware) JwtAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, httputil.MakeResponse(
				http.StatusUnauthorized,
				"Authorization header is required",
				nil,
			))
		}

		token = strings.TrimSpace(strings.Replace(token, "Bearer", "", 1))
		user, err := m.us.Authorize(c.Request().Context(), token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, httputil.MakeResponse(
				http.StatusUnauthorized,
				err.Error(),
				nil,
			))
		}
		// set userId and token on context for later uses
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, domain.TOKEN, token)
		ctx = context.WithValue(ctx, domain.USER, user)
		c.SetRequest(c.Request().Clone(ctx))

		return next(c)
	}
}
