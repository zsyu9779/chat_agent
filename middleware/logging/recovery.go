package logging

import (
	"chat_agent/logger"
	error2 "chat_agent/services/error"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"net/http"
	"net/http/httputil"
)

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				goErr := errors.Wrap(err, 3)
				reset := string([]byte{27, 91, 48, 109})
				logger.Error(fmt.Sprintf(`
[Recovery] panic recovered:
Request : %s
Error :%s
Stack :%s
Reset :%s
`, httprequest, goErr.Error(), goErr.Stack(), reset))
				f(c, err)
			}
		}()
		c.Next() // execute all the handlers
	}
}

func RecoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": error2.ServerErrorCode,
		"msg":  err,
		"body": map[string]interface{}{},
	})
	c.Abort()
}
