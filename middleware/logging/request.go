package logging

import (
	"chat_agent/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func RegisterLog() gin.HandlerFunc {
	logger.Init()
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			logFormat := map[string]interface{}{
				"response_time": params.TimeStamp.Format("2006/01/02 - 15:04:05"),
				"http_code":     params.StatusCode,
				"latency_time":  params.Latency,
				"client_ip":     params.Latency,
				"method":        params.Method,
				"path":          params.Path,
				"error_message": params.ErrorMessage,
			}
			b, _ := json.Marshal(logFormat)
			str := string(b)
			logger.Info(str)
			return str
		},
		Output: WriterLevel(logger.InfoLevel),
		SkipPaths: []string{
			"/ping",
		},
	})
}
