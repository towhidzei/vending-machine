package product

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/algo"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) Buy(ctx context.Context, cart map[uint]uint) (*domain.Bill, error) {
	const op string = "product.service.Buy"

	u, err := domain.UserFromContext(ctx)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	if u.Role != domain.BUYER {
		return nil, domain.ErrPermissionDenied
	}

	s.pl.Lock()
	defer s.pl.Unlock()

	products := make([]domain.Product, 0, len(cart))
	items := make([]domain.Item, 0, len(products))
	var totalPrice uint

	for pid, count := range cart {
		p, err := s.pr.FindById(ctx, pid)
		if err != nil {
			logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
			return nil, domain.ErrProductNotFound
		}
		// check product availability
		if p.Count < count {
			return nil, domain.ErrInsufficientProductsAmount
		}
		// decrease product amount
		p.Count -= count
		products = append(products, *p)
		items = append(items, domain.Item{
			Name:  p.Name,
			Count: count,
			Price: count * p.Price,
		})
		// increase total price
		totalPrice += count * p.Price
		// check user balance
		if u.Deposit < totalPrice {
			return nil, domain.ErrInsufficientBalance
		}
	}

	// passing same tx object in the context
	// when we rollback(commit) one repo,
	// another will rollback(commit) too
	ctx, pr := s.pr.BeginTransaction(ctx)
	ctx, ur := s.ur.BeginTransaction(ctx)

	for _, p := range products {
		err = pr.Update(ctx, &p)
		if err != nil {
			pr.Rollback()
			logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
			return nil, domain.ErrInternalServer
		}
	}

	u.Deposit -= totalPrice
	err = ur.Update(ctx, u)
	if err != nil {
		ur.Rollback()
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}
	// calculating remaining user deposit by valid coins
	refund := algo.MinimumNumberOfElementsWhoseSumIs(domain.Coins, u.Deposit)

	// there is no difference to call pr.Commit()
	// they are in the same transaction
	ur.Commit()
	return &domain.Bill{
		TotalSpent: totalPrice,
		Items:      items,
		Refund:     refund,
	}, nil
}
