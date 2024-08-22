package pgsql

import (
	"github.com/apm-dev/vending-machine/domain"
	"gorm.io/gorm"
)

type Product struct {
	Name     string `gorm:"column:name"`
	Count    uint   `gorm:"column:count"`
	Price    uint   `gorm:"column:cost"`
	SellerID uint   `gorm:"column:seller_id"`
	// gorm model contains id, created_at, updated_at, deleted_at by default
	gorm.Model
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) FromDomain(product domain.Product) {
	p.ID = product.Id
	p.Name = product.Name
	p.Count = product.Count
	p.Price = product.Price
	p.SellerID = product.SellerId
}

func (p *Product) ToDomain() *domain.Product {
	return &domain.Product{
		Id:       p.ID,
		Name:     p.Name,
		Count:    p.Count,
		Price:    p.Price,
		SellerId: p.SellerID,
	}
}
