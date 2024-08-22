package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/domain/mocks"
	"github.com/apm-dev/vending-machine/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Service_Buy(t *testing.T) {
	type args struct {
		ctx  context.Context
		cart map[uint]uint
	}
	type wants struct {
		err  error
		bill *domain.Bill
	}
	type testCase struct {
		name    string
		prepare func()
		args    args
		wants   wants
	}

	pr := new(mocks.ProductRepository)
	ur := new(mocks.UserRepository)
	valueCtx := "*context.valueCtx"
	// cart => productID : count
	normalCart := map[uint]uint{1: 2, 2: 1}
	normalBuyer := &domain.User{
		Id:      1,
		Role:    domain.BUYER,
		Deposit: 500,
	}
	normalContext := context.WithValue(context.Background(), domain.USER, normalBuyer)

	cake := &domain.Product{Id: 1, Name: "Cake", Price: 5, Count: 500}
	soda := &domain.Product{Id: 2, Name: "Soda", Price: 10, Count: 500}

	testCases := []testCase{
		{
			name: "should succeed when a buyer with sufficient balance request valid products",
			prepare: func() {
				pr.On("FindById", mock.Anything, mock.Anything).
					Return(cake, nil).Once()

				pr.On("FindById", mock.Anything, mock.Anything).
					Return(soda, nil).Once()

				pr.On("BeginTransaction", mock.Anything).
					Return(normalContext, pr).Once()
				ur.On("BeginTransaction", mock.Anything).
					Return(normalContext, ur).Once()

				pr.On("Update", mock.Anything, mock.Anything).
					Return(nil).Twice()
				ur.On("Update", mock.Anything, mock.Anything).
					Return(nil).Once()

				ur.On("Commit").Once()
			},
			args: args{
				ctx: context.WithValue(context.Background(), domain.USER, &domain.User{
					Role:    domain.BUYER,
					Deposit: 55,
				}),
				cart: map[uint]uint{1: 2, 2: 1},
			},
			wants: wants{
				err: nil,
				bill: &domain.Bill{
					TotalSpent: 20,
					Items: []domain.Item{
						{Name: "Cake", Count: 2, Price: 10},
						{Name: "Soda", Count: 1, Price: 10},
					},
					Refund: []uint{20, 10, 5},
				},
			},
		},
		{
			name:    "should fail when user is missing from context",
			prepare: func() {},
			args: args{
				ctx:  context.Background(),
				cart: normalCart,
			},
			wants: wants{
				err:  domain.ErrInternalServer,
				bill: nil,
			},
		},
		{
			name:    "should fail when other roles except buyer request",
			prepare: func() {},
			args: args{
				ctx: context.WithValue(context.Background(), domain.USER, &domain.User{
					Role: domain.SELLER,
				}),
				cart: normalCart,
			},
			wants: wants{
				err:  domain.ErrPermissionDenied,
				bill: nil,
			},
		},
		{
			name: "should fail when product not found",
			prepare: func() {
				pr.On("FindById", mock.AnythingOfType(valueCtx), mock.Anything).
					Return(nil, domain.ErrProductNotFound).Once()
			},
			args: args{
				ctx:  normalContext,
				cart: normalCart,
			},
			wants: wants{
				err:  domain.ErrProductNotFound,
				bill: nil,
			},
		},
		{
			name: "should fail when requested product has no sufficient amount",
			prepare: func() {
				pr.On("FindById", mock.AnythingOfType(valueCtx), mock.Anything).
					Return(&domain.Product{Price: 10, Count: 2}, nil).Once()
			},
			args: args{
				ctx:  normalContext,
				cart: map[uint]uint{1: 5},
			},
			wants: wants{
				err:  domain.ErrInsufficientProductsAmount,
				bill: nil,
			},
		},
		{
			name: "should fail when buyer has no sufficient balance",
			prepare: func() {
				pr.On("FindById", mock.AnythingOfType(valueCtx), mock.Anything).
					Return(&domain.Product{Price: 30, Count: 20}, nil).Once()
			},
			args: args{
				ctx: context.WithValue(context.Background(), domain.USER, &domain.User{
					Role:    domain.BUYER,
					Deposit: 15,
				}),
				cart: map[uint]uint{1: 10},
			},
			wants: wants{
				err:  domain.ErrInsufficientBalance,
				bill: nil,
			},
		},
		{
			name: "should fail and rollback changes when products update fail",
			prepare: func() {
				pr.On("FindById", mock.AnythingOfType("*context.valueCtx"), mock.Anything).
					Return(cake, nil).Once()

				pr.On("BeginTransaction", mock.Anything).
					Return(normalContext, pr).Once()
				ur.On("BeginTransaction", mock.Anything).
					Return(normalContext, ur).Once()

				pr.On("Update", mock.AnythingOfType("*context.valueCtx"), mock.Anything).
					Return(errors.New("failed to update product")).Once()

				pr.On("Rollback").Once()
			},
			args: args{
				ctx:  normalContext,
				cart: map[uint]uint{1: 2},
			},
			wants: wants{
				err:  domain.ErrInternalServer,
				bill: nil,
			},
		},
		{
			name: "should fail and rollback changes when user update fail",
			prepare: func() {
				pr.On("FindById", mock.Anything, mock.Anything).
					Return(cake, nil).Once()

				pr.On("BeginTransaction", mock.Anything).
					Return(normalContext, pr).Once()
				ur.On("BeginTransaction", mock.Anything).
					Return(normalContext, ur).Once()

				pr.On("Update", mock.Anything, mock.Anything).
					Return(nil).Once()
				ur.On("Update", mock.Anything, mock.Anything).
					Return(errors.New("failed to update user")).Once()

				ur.On("Rollback").Once()
			},
			args: args{
				ctx:  normalContext,
				cart: map[uint]uint{1: 2},
			},
			wants: wants{
				err:  domain.ErrInternalServer,
				bill: nil,
			},
		},
	}

	svc := product.InitService(pr, ur)

	for _, tc := range testCases {
		// arrange
		tc.prepare()
		// action
		bill, err := svc.Buy(tc.args.ctx, tc.args.cart)
		// assert
		if tc.wants.err != nil {
			assert.ErrorIs(t, err, tc.wants.err, tc.name)
			assert.EqualValues(t, tc.wants.bill, bill, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
			assert.EqualValues(t, tc.wants.bill, bill, tc.name)
		}
	}
	pr.AssertExpectations(t)
	ur.AssertExpectations(t)
}
