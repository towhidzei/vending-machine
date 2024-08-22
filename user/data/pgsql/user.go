package pgsql

import (
	"github.com/apm-dev/vending-machine/domain"
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"uniqueIndex;size:32;column:username"`
	Password string `gorm:"size:256;column:password"`
	Role     string `gorm:"size:32;column:role"`
	Deposit  uint   `gorm:"column:deposit"`
	gorm.Model
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) FromDomain(user *domain.User) {
	u.ID = user.Id
	u.Username = user.Username
	u.Password = user.Password
	u.Role = string(user.Role)
	u.Deposit = user.Deposit
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		Id:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Role:      domain.Role(u.Role),
		Deposit:   u.Deposit,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt.Time,
	}
}
