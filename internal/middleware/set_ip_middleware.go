package middleware

import (
	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/ctxutil"
	"github.com/gin-gonic/gin"
)

func SetIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ctxutil.SetValueToContext(c.Request.Context(), ctxutil.ContextValue{
			Key:   constant.IPAddressKey,
			Value: c.ClientIP(),
		})
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
