package repository

import (
	"context"
	"time"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/handling"
	"gorm.io/gorm"
)

type QueryRepository interface {
	ExecuteRaw(
		ctx context.Context,
		sql string,
		values ...any,
	) ([]map[string]any, error)
}

type queryRepository struct {
	db *gorm.DB
}

func (r *queryRepository) ExecuteRaw(ctx context.Context, sql string, values ...any) ([]map[string]any, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)
	defer cancel()

	raws := []map[string]any{}

	query := r.db.WithContext(ctx).Raw(sql, values...).Scan(&raws)

	if query.Error != nil {
		return raws, handling.NewErrorWrapper(constant.CodeDBError, "Error query raw", nil, query.Error)
	}

	return raws, nil
}
