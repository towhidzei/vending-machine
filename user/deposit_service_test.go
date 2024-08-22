package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/domain/mocks"
	"github.com/apm-dev/vending-machine/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Service_Deposit(t *testing.T) {
	type args struct {
		ctx  context.Context
		coin domain.Coin
	}
	type wants struct {
		err     error
		balance uint
	}
	type testCase struct {
		name    string
		timeout time.Duration
		prepare func()
		args    args
		wants   wants
	}

	tout := time.Second * 2
	timerCtx := "*context.timerCtx"
	ur := new(mocks.UserRepository)

	testCases := []testCase{
		{
			name:    "should succeed when buyer deposits valid coin",
			timeout: tout,
			prepare: func() {
				u := &domain.User{
					Id:      1,
					Role:    domain.BUYER,
					Deposit: 0,
				}
				ur.On("FindById",
					mock.AnythingOfType(timerCtx), uint(1),
				).Return(u, nil).Once()
				ur.On("Update",
					mock.AnythingOfType(timerCtx),
					mock.AnythingOfType("*domain.User"),
				).Return(nil).Once()
			},
			args: args{
				ctx: context.WithValue(context.Background(), domain.USER, &domain.User{
					Id:      1,
					Role:    domain.BUYER,
					Deposit: 0,
				}),
				coin: 50,
			},
			wants: wants{
				err:     nil,
				balance: 50,
			},
		},
		{
			name:    "should fail when buyer deposits invalid coin",
			timeout: tout,
			prepare: func() {},
			args: args{
				ctx:  context.WithValue(context.Background(), domain.USER, &domain.User{}),
				coin: 31,
			},
			wants: wants{
				err:     domain.ErrInvalidCoin,
				balance: 0,
			},
		},
		{
			name:    "should fail when seller deposits",
			timeout: tout,
			prepare: func() {
				u := &domain.User{
					Id:      1,
					Role:    domain.SELLER,
					Deposit: 0,
				}
				ur.On("FindById",
					mock.AnythingOfType(timerCtx), uint(1),
				).Return(u, nil).Once()
			},
			args: args{
				ctx: context.WithValue(context.Background(), domain.USER, &domain.User{
					Id:      1,
					Role:    domain.SELLER,
					Deposit: 0,
				}),
				coin: 5,
			},
			wants: wants{
				err:     domain.ErrPermissionDenied,
				balance: 0,
			},
		},
		{
			name:    "should fail when user is missing from context",
			timeout: tout,
			prepare: func() {},
			args: args{
				ctx:  context.Background(),
				coin: 5,
			},
			wants: wants{
				err:     domain.ErrInternalServer,
				balance: 0,
			},
		},
		{
			name:    "should fail when user repository returns error",
			timeout: tout,
			prepare: func() {
				ur.On("FindById",
					mock.AnythingOfType(timerCtx), uint(1),
				).Return(nil, domain.ErrUserNotFound).Once()
			},
			args: args{
				ctx:  context.WithValue(context.Background(), domain.USER, &domain.User{Id: 1, Role: domain.BUYER}),
				coin: 5,
			},
			wants: wants{
				err:     domain.ErrInternalServer,
				balance: 0,
			},
		},
	}

	for _, tc := range testCases {
		// arrange
		tc.prepare()
		svc := user.InitService(ur, nil, nil, tc.timeout)
		// action
		balance, err := svc.Deposit(tc.args.ctx, tc.args.coin)
		// assert
		if tc.wants.err != nil {
			assert.ErrorIs(t, err, tc.wants.err, tc.name)
			assert.EqualValues(t, tc.wants.balance, balance, tc.name)
		} else {
			assert.NoError(t, err, tc.name)
			assert.EqualValues(t, tc.wants.balance, balance, tc.name)
		}
	}
	ur.AssertExpectations(t)
}
