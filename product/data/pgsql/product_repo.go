package pgsql

import (
	"context"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/pgsqlhelper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func InitProductRepository(db *gorm.DB) domain.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) BeginTransaction(ctx context.Context) (context.Context, domain.ProductRepository) {
	if tx, ok := pgsqlhelper.TransactionFromContext(ctx); ok {
		return ctx, InitProductRepository(tx)
	}
	tx := r.db.Begin()
	ctx = pgsqlhelper.TransactionToContext(ctx, tx)
	return ctx, InitProductRepository(tx)
}

func (r *ProductRepository) Commit() {
	r.db.Commit()
}

func (r *ProductRepository) Rollback() {
	r.db.Rollback()
}

func (r *ProductRepository) Insert(ctx context.Context, p domain.Product) (uint, error) {
	const op string = "product.data.pgsql.product_repo.Insert"

	dbp := new(Product)
	dbp.FromDomain(p)

	err := r.db.WithContext(ctx).Create(&dbp).Error
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	return dbp.ID, nil
}

func (r *ProductRepository) FindById(ctx context.Context, id uint) (*domain.Product, error) {
	const op string = "product.data.pgsql.product_repo.FindById"

	dbp := new(Product)

	err := r.db.WithContext(ctx).First(&dbp, id).Error
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return dbp.ToDomain(), nil
}

func (r *ProductRepository) List(ctx context.Context) ([]domain.Product, error) {
	const op string = "product.data.pgsql.product_repo.List"

	var dbps []Product

	err := r.db.WithContext(ctx).Find(&dbps).Error
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var ps = make([]domain.Product, len(dbps))
	for i, dbp := range dbps {
		ps[i] = *dbp.ToDomain()
	}

	return ps, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	const op string = "product.data.pgsql.product_repo.Update"

	dbp := new(Product)
	dbp.FromDomain(*p)

	err := r.db.WithContext(ctx).Save(&dbp).Error
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	const op string = "product.data.pgsql.product_repo.Delete"

	err := r.db.WithContext(ctx).Delete(&Product{}, id).Error
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
