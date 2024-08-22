package rest

import (
	"net/http"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/product/presentation/rest/requests"
	"github.com/labstack/echo"
)

func (p *ProductHandler) Buy(c echo.Context) error {
	req := new(requests.Buy)
	err := httputil.BindAndValidate(c, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, err.Error(), nil,
		))
	}
	if len(req.Cart) == 0 {
		return c.JSON(http.StatusBadRequest, httputil.MakeResponse(
			http.StatusBadRequest, domain.ErrInvalidParams.Error(), nil,
		))
	}

	bill, err := p.ps.Buy(c.Request().Context(), req.Cart)

	return checkErrorThenResponse(c, err, bill)
}
