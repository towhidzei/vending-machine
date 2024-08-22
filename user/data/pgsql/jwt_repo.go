package pgsql

import (
	"context"
	"time"

	"github.com/apm-dev/vending-machine/domain"
	"github.com/apm-dev/vending-machine/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type JwtRepository struct {
	db *gorm.DB
}

func InitJwtRepository(db *gorm.DB) domain.JwtRepository {
	return &JwtRepository{db}
}

func (r *JwtRepository) Insert(ctx context.Context, userId uint, token string, ttl time.Duration) error {
	const op string = "user.data.pgsql.jwt_repo.Insert"

	jwt := &JWT{
		Token:     token,
		UserID:    userId,
		ExpiredAt: time.Now().Add(ttl),
	}

	err := r.db.WithContext(ctx).Create(jwt).Error
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return err
	}

	return nil
}

func (r *JwtRepository) Exists(ctx context.Context, token string) (bool, error) {
	const op string = "user.data.pgsql.jwt_repo.Exists"

	result := r.db.WithContext(ctx).First(&JWT{}, "token = ?", token)
	if result.Error != nil || result.RowsAffected == 0 {
		return false, errors.Wrap(result.Error, op)
	}
	return true, nil
}

func (r *JwtRepository) UserTokensCount(ctx context.Context, uid uint) (uint, error) {
	const op string = "user.data.pgsql.jwt_repo.UserTokensCount"

	var count int64

	err := r.db.WithContext(ctx).Model(&JWT{}).Where("user_id = ?", uid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	return uint(count), nil
}

func (r *JwtRepository) DeleteTokensOfUserExcept(ctx context.Context, userId uint, exceptionToken string) error {
	const op string = "user.data.pgsql.jwt_repo.DeleteTokensOfUserExcept"

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND token != ?", userId, exceptionToken).
		Delete(&JWT{}).Error

	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
