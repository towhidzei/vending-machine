package user

import (
	"context"
	"fmt"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
)

// Register creates new user and return jwt token or error
func (s *Service) Register(ctx context.Context, uname, pass string, role domain.Role) (string, error) {
	const op string = "user.service.Register"
	// create domain user object
	user, err := domain.NewUser(uname, pass, role)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}
	// persist user
	user.Id, err = s.ur.Insert(ctx, *user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return "", domain.ErrUserAlreadyExists
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	// generate jwt token with user claims
	token, err := s.jwt.Generate(user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	// persist jwt token
	err = s.jr.Insert(ctx, user.Id, token, s.jwt.tokenDuration)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	logger.Log(logger.INFO, fmt.Sprintf("%s registered", uname))

	return token, nil
}

// Login checks credentials, generate and return jwt token and a boolean
// which says is there another active session using this account or not
func (s *Service) Login(ctx context.Context, uname, pass string) (string, bool, error) {
	const op string = "user.service.Login"
	// fetch user from db
	user, err := s.ur.FindByUsername(ctx, uname)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", false, domain.ErrUserNotFound
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", false, domain.ErrInternalServer
	}
	// check password
	if !user.CheckPassword(pass) {
		logger.Log(logger.INFO, errors.Wrap(domain.ErrWrongCredentials, user.Username).Error())
		return "", false, domain.ErrWrongCredentials
	}
	// generate jwt token with user claims
	token, err := s.jwt.Generate(user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", false, domain.ErrInternalServer
	}

	// persist jwt token
	err = s.jr.Insert(ctx, user.Id, token, s.jwt.tokenDuration)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", false, domain.ErrInternalServer
	}

	// check is there any other active session or not
	count, err := s.jr.UserTokensCount(ctx, user.Id)
	if err != nil {
		logger.Log(logger.WARN, errors.Wrap(err, op).Error())
	}

	logger.Log(logger.INFO, fmt.Sprintf(
		"%s logged-in", user.Username,
	))

	return token, count > 1, nil
}

// Authorize parses jwt token and return related user
func (s *Service) Authorize(ctx context.Context, token string) (*domain.User, error) {
	const op string = "user.service.Authorize"

	// verify and get claims of token
	claims, err := s.jwt.Verify(token)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			logger.Log(logger.INFO, errors.Wrap(err, op).Error())
			return nil, domain.ErrInvalidToken
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	// check token existans in db to see is it still active or not
	exists, err := s.jr.Exists(ctx, token)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	if !exists {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.ur.FindById(ctx, claims.Id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

// TerminateActiveSessions terminates all other active sessions
func (s *Service) TerminateActiveSessions(ctx context.Context) error {
	const op string = "user.service.TerminateActiveSessions"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.WARN, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	token, err := domain.TokenFromContext(ctx)
	if err != nil {
		logger.Log(logger.WARN, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	err = s.jr.DeleteTokensOfUserExcept(ctx, u.Id, token)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	return nil
}
