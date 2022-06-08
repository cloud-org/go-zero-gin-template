package middleware

import (
	"bytes"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/core/utils"
)

// LogxLogger go-zero gin logx logger
func LogxLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := utils.NewElapsedTimer()
		c.Next()
		duration := timer.Duration()

		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf("%s - %s - %s - %s - %d - %s - %s",
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			timex.ReprOfDuration(duration),
			c.Request.UserAgent(),
		))

		logx.Info(buf.String())
	}
}
