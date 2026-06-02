package route

import (
	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/handling"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB, cfg *constant.Config) {

	router.GET("/echo", func(c *gin.Context) {
		res := handling.ResponseSuccess("test", "Echo test", "")
		c.JSON(res.Code, &res)
	})

}
