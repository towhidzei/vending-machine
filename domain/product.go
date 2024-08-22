package domain

import "context"

type Product struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	Count    uint   `json:"count"`
	Price    uint   `json:"price"`
	SellerId uint   `json:"seller_id"`
}

type Bill struct {
	TotalSpent uint   `json:"total_spent"`
	Items      []Item `json:"items"`
	Refund     []uint `json:"refund"`
}

type Item struct {
	Name  string `json:"name"`
	Count uint   `json:"count"`
	Price uint   `json:"price"`
}

func NewProduct(name string, amountAvailable, cost, sellerId uint) *Product {
	return &Product{
		Name:     name,
		Count:    amountAvailable,
		Price:    cost,
		SellerId: sellerId,
	}
}

type ProductService interface {
	Add(ctx context.Context, name string, amount, cost uint) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, id uint, name string, amount, cost uint) (*Product, error)
	Delete(ctx context.Context, id uint) error
	Buy(ctx context.Context, cart map[uint]uint) (*Bill, error)
}

type ProductRepository interface {
	DBTransaction
	BeginTransaction(ctx context.Context) (context.Context, ProductRepository)
	Insert(ctx context.Context, p Product) (uint, error)
	FindById(ctx context.Context, id uint) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uint) error
}
