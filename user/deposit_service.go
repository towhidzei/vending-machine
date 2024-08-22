package user

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/algo"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
)

// Deposit increases buyer(user) deposit
func (s *Service) Deposit(ctx context.Context, coin domain.Coin) (uint, error) {
	const op string = "user.service.Deposit"

	ctx, cancel := context.WithTimeout(ctx, s.dtout)
	defer cancel()

	if !coin.IsValid() {
		return 0, domain.ErrInvalidCoin
	}

	user, err := s.refetchContextUserFromDB(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return 0, domain.ErrInternalServer
	}

	if user.Role != domain.BUYER {
		return 0, domain.ErrPermissionDenied
	}

	// use locks because of concurrent requests
	// we make sure to add user deposits without conflict
	// there are better ways to handle this kind of issue
	// but it's just for POC
	s.dl.Lock()
	defer s.dl.Unlock()

	user.AddDeposit(coin)

	err = s.ur.Update(ctx, user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return 0, domain.ErrInternalServer
	}

	return user.Deposit, nil
}

// ResetDeposit reset buyer(user) deposits back to zero
func (s *Service) ResetDeposit(ctx context.Context) ([]uint, error) {
	const op string = "user.service.ResetDeposit"

	ctx, cancel := context.WithTimeout(ctx, s.dtout)
	defer cancel()

	user, err := s.refetchContextUserFromDB(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	if user.Role != domain.BUYER {
		return nil, domain.ErrPermissionDenied
	}

	s.dl.Lock()
	defer s.dl.Unlock()

	deposit := user.Deposit
	user.ResetDeposit()

	err = s.ur.Update(ctx, user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// calculate user refund with valid coins
	refund := algo.MinimumNumberOfElementsWhoseSumIs(domain.Coins, deposit)
	return refund, nil
}
