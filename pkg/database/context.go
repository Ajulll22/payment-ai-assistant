package database

import (
	"context"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"gorm.io/gorm"
)

func GetTxFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(constant.DBTransactionKey).(*gorm.DB)
	return tx, ok
}

func SetTxToContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, constant.DBTransactionKey, tx)
}
