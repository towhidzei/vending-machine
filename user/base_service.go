package user

import (
	"context"
	"sync"
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/algo"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
)

type Service struct {
	ur  domain.UserRepository
	jr  domain.JwtRepository
	jwt *JWTManager
	// deposit timeout
	dtout time.Duration
	// deposit lock
	dl sync.RWMutex
}

var UserService *Service

func InitService(
	ur domain.UserRepository,
	jr domain.JwtRepository,
	jwt *JWTManager,
	dtout time.Duration,
) domain.UserService {
	if UserService == nil {
		UserService = &Service{
			ur: ur, jr: jr, jwt: jwt,
			dtout: dtout,
		}
	}
	return UserService
}

func (s *Service) Update(ctx context.Context, passwd string) error {
	const op string = "user.service.Update"

	user, err := s.refetchContextUserFromDB(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	err = user.SetPassword(passwd)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	err = s.ur.Update(ctx, user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}
	return nil
}

func (s *Service) Delete(ctx context.Context) ([]uint, error) {
	const op string = "user.service.Delete"

	user, err := s.refetchContextUserFromDB(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	err = s.ur.Delete(ctx, user.Id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	refund := algo.MinimumNumberOfElementsWhoseSumIs(domain.Coins, user.Deposit)
	return refund, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*domain.User, error) {
	const op string = "user.service.Get"

	user, err := s.refetchContextUserFromDB(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	// all users except ADMIN ones only can fetch themselves profile
	if user.Role != domain.ADMIN {
		return user, nil
	}

	user, err = s.ur.FindById(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (s *Service) List(ctx context.Context) ([]domain.User, error) {
	const op string = "user.service.List"

	user, err := s.refetchContextUserFromDB(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	// only ADMINs can fetch list of users
	if user.Role != domain.ADMIN {
		return nil, domain.ErrPermissionDenied
	}

	users, err := s.ur.List(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrUserNotFound
	}

	return users, nil
}
