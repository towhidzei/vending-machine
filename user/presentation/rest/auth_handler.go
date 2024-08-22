package rest

import (
	"net/http"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/user/presentation/rest/requests"
	"github.com/labstack/echo"
)

func (h *UserHandler) Register(c echo.Context) error {
	req := new(requests.Register)
	if err := httputil.BindAndValidate(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	token, err := h.us.Register(
		c.Request().Context(),
		req.Username, req.Password,
		domain.Role(req.Role),
	)
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "welcome "+req.Username, echo.Map{"token": token},
	))
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(requests.Login)
	if err := httputil.BindAndValidate(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	token, activeSessions, err := h.us.Login(
		c.Request().Context(),
		req.Username, req.Password,
	)
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}
	// notify the user if there was another active session
	msg := "welcome " + req.Username
	if activeSessions {
		msg += ", there is another active session."
	}
	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, msg, echo.Map{"token": token},
	))
}

func (h *UserHandler) LogoutAll(c echo.Context) error {
	err := h.us.TerminateActiveSessions(c.Request().Context())
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK,
		"All other active sessions have been terminated.",
		nil,
	))
}

func (h *UserHandler) UpdatePassword(c echo.Context) error {
	req := new(requests.UpdatePassword)
	if err := httputil.BindAndValidate(c, req); err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	err := h.us.Update(c.Request().Context(), req.Password)
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "Password updated.", nil,
	))
}

func (h *UserHandler) DeleteAccount(c echo.Context) error {
	refund, err := h.us.Delete(c.Request().Context())
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "Account deleted.", refund,
	))
}
