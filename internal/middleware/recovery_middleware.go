package middleware

import (
	"runtime/debug"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/database"
	"github.com/Ajulll22/payment-ai-assistant/pkg/handling"
	"github.com/Ajulll22/payment-ai-assistant/pkg/logger"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware(cfg constant.LogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				ctx := c.Request.Context()
				log := logger.FromContext(ctx, cfg)

				if tx, ok := database.GetTxFromContext(ctx); ok {
					// rollback transaksi
					_ = tx.Rollback().Error
				}

				log.App.Error("-------------------------------------------------------------")
				log.App.Error("Panic recovered:", r)
				log.App.Error(string(debug.Stack()))
				log.App.Error("-------------------------------------------------------------")

				newToken := c.GetString("new_token")

				res := handling.ResponseError(ctx, cfg, handling.NewErrorWrapper(constant.CodeInternalServer, "Error recovery", nil, nil), newToken)

				c.JSON(res.Code, &res)
			}
		}()
		c.Next()
	}
}
