package utils

import "github.com/gin-gonic/gin"

type StatusCode int

const (
	OK StatusCode = -iota
	InvalidBody
	InvalidID
	ExecutionError
	RuntimeError
)

var ErrorMsgMap = map[StatusCode]string{
	InvalidBody:    "Invalid body",
	InvalidID:      "Invalid ID",
	ExecutionError: "Execution error",
	RuntimeError:   "Runtime error",
}

func Respond(c *gin.Context, code StatusCode, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"data": data,
	})
}

func RespondError(c *gin.Context, code StatusCode, err error) {
	c.JSON(400, gin.H{
		"code":  code,
		"msg":   ErrorMsgMap[code],
		"error": err,
	})
}

func RespondSuccess(c *gin.Context, data interface{}) {
	Respond(c, OK, data)
}
