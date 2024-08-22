package rest

import (
	"net/http"
	"strconv"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/product/presentation/rest/requests"
	"github.com/labstack/echo"
)

type ProductHandler struct {
	ps domain.ProductService
}

// InitUserHandler
// e echo instance to define normal routes (no authorization need)
// auth echo group which uses auth middleware
func InitProductHandler(e *echo.Echo, auth *echo.Group, ps domain.ProductService) *ProductHandler {
	h := &ProductHandler{ps: ps}

	pg := e.Group("/products")
	pg.GET("/", h.List)

	pg = auth.Group("/products")
	pg.POST("/", h.Add)
	pg.PUT("/:id", h.Update)
	pg.DELETE("/:id", h.Delete)

	pg.POST("/buy", h.Buy)

	return h
}

func (h *ProductHandler) List(c echo.Context) error {
	ps, err := h.ps.List(c.Request().Context())
	return checkErrorThenResponse(c, err, ps)
}

func (h *ProductHandler) Add(c echo.Context) error {
	req := new(requests.AddProduct)
	err := httputil.BindAndValidate(c, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	p, err := h.ps.Add(c.Request().Context(),
		req.Name, req.Count, req.Price,
	)

	return checkErrorThenResponse(c, err, p)
}

func (h *ProductHandler) Update(c echo.Context) error {
	req := new(requests.AddProduct)
	err := httputil.BindAndValidate(c, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	p, err := h.ps.Update(c.Request().Context(),
		uint(id), req.Name, req.Count, req.Price,
	)

	return checkErrorThenResponse(c, err, p)
}

func (p *ProductHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}

	err = p.ps.Delete(c.Request().Context(), uint(id))

	return checkErrorThenResponse(c, err, p)
}

func checkErrorThenResponse(c echo.Context, err error, content interface{}) error {
	if err != nil {
		status := httputil.StatusCode(err)
		return c.JSON(status, httputil.MakeResponse(
			status, err.Error(), nil,
		))
	}

	return c.JSON(http.StatusOK, httputil.MakeResponse(
		http.StatusOK, "", content,
	))
}
