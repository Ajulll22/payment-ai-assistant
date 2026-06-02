package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/Ajulll22/payment-ai-assistant/pkg/constant"
	"github.com/Ajulll22/payment-ai-assistant/pkg/handling"
	"github.com/Ajulll22/payment-ai-assistant/pkg/logger"
	"github.com/Ajulll22/payment-ai-assistant/pkg/security"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(cfg constant.LogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		log := logger.FromContext(c.Request.Context(), cfg)
		defer log.Close()

		ctx := logger.InjectIntoContext(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)

		// check config for write request
		if !cfg.RequestStatus {
			return
		}

		// construct request from body
		body, err := io.ReadAll(c.Request.Body)

		if err != nil {

			res := handling.ResponseError(ctx, cfg, err, "")

			c.JSON(res.Code, &res)

			return
		}

		requestBody := make(map[string]any)
		json.Unmarshal(body, &requestBody)

		requestBodyMasking := security.MaskSensitive(requestBody)

		log.App.Process("Log start ===========================")
		log.App.Process("URL :", c.Request.URL)
		log.App.Process("Method :", c.Request.Method)
		log.App.Process("IP Address :", c.ClientIP())
		log.App.Process("Request Data ===========================")
		log.App.Map(requestBodyMasking)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Pass on to the next-in-chain
		c.Next()

		// construct log response
		log.App.Process("Response Data ===========================")
		log.App.Process("Status :", c.Writer.Status())

		responseString := blw.body.String()
		b := []byte(responseString)
		responsetBody := make(map[string]any)
		json.Unmarshal(b, &responsetBody)

		responseBodyMasking := security.MaskSensitive(responsetBody)

		log.App.Map(responseBodyMasking)
		log.App.Process("Log end ===========================")

	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
