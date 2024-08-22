package product

import (
	"context"
	"sync"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
)

type Service struct {
	pr domain.ProductRepository
	ur domain.UserRepository
	pl sync.RWMutex
}

func InitService(pr domain.ProductRepository, ur domain.UserRepository) domain.ProductService {
	return &Service{pr: pr, ur: ur}
}

func (s *Service) Add(ctx context.Context, name string, amount uint, cost uint) (*domain.Product, error) {
	const op string = "product.service.Add"

	if cost%5 != 0 {
		return nil, domain.ErrInvalidCost
	}
	cu, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrUserNotFound
	}

	if cu.Role != domain.SELLER {
		return nil, domain.ErrPermissionDenied
	}

	p := domain.NewProduct(name, amount, cost, cu.Id)

	p.Id, err = s.pr.Insert(ctx, *p)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return p, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Product, error) {
	const op string = "product.service.List"

	s.pl.RLock()
	defer s.pl.RUnlock()

	ps, err := s.pr.List(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return ps, nil
}

func (s *Service) Update(ctx context.Context, id uint, name string, amount, cost uint) (*domain.Product, error) {
	const op string = "product.service.Update"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// only sellers can update products
	if u.Role != domain.SELLER {
		return nil, domain.ErrPermissionDenied
	}

	s.pl.Lock()
	defer s.pl.Unlock()

	p, err := s.pr.FindById(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// only related seller can update it
	if p.SellerId != u.Id {
		return nil, domain.ErrPermissionDenied
	}

	p.Name = name
	p.Count = amount
	p.Price = cost

	err = s.pr.Update(ctx, p)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return p, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	const op string = "product.service.Delete"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrUserNotFound
	}

	if u.Role != domain.SELLER {
		return domain.ErrPermissionDenied
	}

	p, err := s.pr.FindById(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	// only related seller can delete it
	if p.SellerId != u.Id {
		return domain.ErrPermissionDenied
	}

	s.pl.Lock()
	defer s.pl.Unlock()

	err = s.pr.Delete(ctx, id)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return domain.ErrInternalServer
	}

	return nil
}
