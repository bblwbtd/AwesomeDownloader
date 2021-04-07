package utils

import "github.com/gin-gonic/gin"

type StatusCode int

const (
	OK StatusCode = -iota
	InvalidBody
	ExecutionError
)

func Respond(c *gin.Context, code StatusCode, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"data": data,
	})
}

func RespondError(c *gin.Context, code StatusCode, msg string) {
	c.JSON(400, gin.H{
		"code": code,
		"msg":  msg,
	})
}

func RespondSuccess(c *gin.Context, data interface{}) {
	Respond(c, OK, data)
}
