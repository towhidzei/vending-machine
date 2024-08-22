package pgsql

import (
	"time"
)

type JWT struct {
	Token     string `gorm:"primaryKey;column:token"`
	UserID    uint   `gorm:"column:user_id"`
	User      User
	ExpiredAt time.Time `gorm:"column:expired_at"`
}

func (j *JWT) TableName() string {
	return "jwt_tokens"
}
