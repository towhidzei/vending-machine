package pgsqlhelper

import (
	"context"

	"gorm.io/gorm"
)

type txContextKey string

const txn txContextKey = "TransactionContextKey"

func TransactionFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txn).(*gorm.DB)
	return tx, ok
}

func TransactionToContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txn, tx)
}
