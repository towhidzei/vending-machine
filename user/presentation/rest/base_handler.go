package rest

import (
	"net/http"
	"strconv"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/labstack/echo"
)

type UserHandler struct {
	us domain.UserService
}

// InitUserHandler
// e echo instance to define normal routes (no authorization need)
// auth echo group which uses auth middleware
func InitUserHandler(e *echo.Echo, auth *echo.Group, us domain.UserService) *UserHandler {
	h := &UserHandler{us}
	// public routes
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	// auth required routes
	auth.POST("/logout/all", h.LogoutAll)
	auth.POST("/deposit", h.Deposit)
	auth.POST("/reset", h.ResetDeposit)
	// user CRUD
	// Restful standard
	// GET 			/users, /users/:id
	// POST 		/users
	// PUT, PATCH 	/users/:id
	// DELETE 		/users/:id
	u := auth.Group("/users")
	u.GET("/", h.List)
	u.GET("/:id", h.Profile)
	u.PATCH("/:id", h.UpdatePassword)
	u.DELETE("/:id", h.DeleteAccount)

	return h
}

func (h *UserHandler) List(c echo.Context) error {
	users, err := h.us.List(c.Request().Context())
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "", users,
	))
}

func (h *UserHandler) Profile(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, ":id must be positive number", nil,
		))
	}

	user, err := h.us.Get(c.Request().Context(), uint(id))
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "", user,
	))
}
