package user

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/pkg/errors"
)

func (s *Service) refetchContextUserFromDB(ctx context.Context) (*domain.User, error) {
	const op string = "user.helper.fetchContextUser"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	user, err := s.ur.FindById(ctx, u.Id)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return user, nil
}
