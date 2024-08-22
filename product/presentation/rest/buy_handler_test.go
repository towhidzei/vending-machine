package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/domain/mocks"
	"github.com/apm-dev/vending-machine/pkg/httputil"
	"github.com/apm-dev/vending-machine/product/presentation/rest"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductHandler_Buy(t *testing.T) {
	type args struct {
		// map of product id => count
		Cart map[uint]uint `json:"cart"`
	}
	type wants struct {
		status int
		bill   *domain.Bill
	}
	type testCase struct {
		name    string
		prepare func()
		args    args  // http request body
		wants   wants // http response body
	}

	ps := new(mocks.ProductService)
	mockBill := &domain.Bill{
		TotalSpent: 20,
		Items: []domain.Item{
			{Name: "P1", Count: 1, Price: 10},
			{Name: "P2", Count: 2, Price: 5},
		},
		Refund: []uint{10, 5},
	}

	testCases := []testCase{
		{
			name: "200 OK and bill",
			prepare: func() {
				ps.On("Buy", mock.Anything, mock.Anything).
					Return(mockBill, nil).Once()
			},
			args: args{map[uint]uint{1: 1, 2: 2}},
			wants: wants{
				status: http.StatusOK,
				bill:   mockBill,
			},
		},
		{
			name:    "400 BadRequest",
			prepare: func() {},
			args: args{
				Cart: nil,
			},
			wants: wants{
				status: http.StatusBadRequest,
				bill:   nil,
			},
		},
		{
			name: "500 InternalServerError",
			prepare: func() {
				ps.On("Buy", mock.Anything, mock.Anything).
					Return(nil, domain.ErrInternalServer).Once()
			},
			args: args{map[uint]uint{1: 2}},
			wants: wants{
				status: http.StatusInternalServerError,
				bill:   nil,
			},
		},
		{
			name: "401 Unauthorized",
			prepare: func() {
				ps.On("Buy", mock.Anything, mock.Anything).
					Return(nil, domain.ErrUnauthorized).Once()
			},
			args: args{map[uint]uint{1: 2}},
			wants: wants{
				status: http.StatusUnauthorized,
				bill:   nil,
			},
		},
		{
			name: "403 Forbidden",
			prepare: func() {
				ps.On("Buy", mock.Anything, mock.Anything).
					Return(nil, domain.ErrPermissionDenied).Once()
			},
			args: args{map[uint]uint{1: 2}},
			wants: wants{
				status: http.StatusForbidden,
				bill:   nil,
			},
		},
	}

	e := echo.New()
	e.Validator = httputil.InitCustomValidator()

	for _, tc := range testCases {
		// arrange
		tc.prepare()

		body, err := json.Marshal(tc.args)
		require.NoError(t, err, tc.name)
		req, err := http.NewRequest(echo.POST, "/products/buy", strings.NewReader(string(body)))
		require.NoError(t, err, tc.name)
		req.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()

		c := e.NewContext(req, response)
		handler := rest.InitProductHandler(e, e.Group(""), ps)
		// action
		err = handler.Buy(c)
		// assert
		assert.NoError(t, err, tc.name)
		assert.Equal(t, tc.wants.status, response.Code, tc.name)
		// checking for expected rest-api response
		if tc.wants.bill != nil {
			type body struct {
				Code    int          `json:"code"`
				Message string       `json:"message"`
				Content *domain.Bill `json:"content"`
			}
			b := new(body)
			err = json.NewDecoder(response.Body).Decode(&b)
			assert.NoError(t, err, tc.name)
			assert.EqualValues(t, tc.wants.bill, b.Content, tc.name)
		}
	}
}
