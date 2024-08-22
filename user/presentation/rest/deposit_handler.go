package rest

import (
	"net/http"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/user/presentation/rest/requests"
	"github.com/labstack/echo"
)

func (h *UserHandler) Deposit(c echo.Context) error {
	req := new(requests.Deposit)
	if err := httputil.BindAndValidate(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	b, err := h.us.Deposit(c.Request().Context(), domain.Coin(req.Coin))
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "", echo.Map{"balance": b},
	))
}

func (h *UserHandler) ResetDeposit(c echo.Context) error {
	refund, err := h.us.ResetDeposit(c.Request().Context())
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "", echo.Map{"refund": refund},
	))
}
