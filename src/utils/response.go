package utils

import "github.com/gin-gonic/gin"

type StatusCode int

const (
	OK StatusCode = -iota
	InvalidBody
	InvalidID
	ExecutionError
)

var ErrorMsgMap = map[StatusCode]string{
	InvalidBody:    "Invalid Body",
	InvalidID:      "Invalid ID",
	ExecutionError: "Execution Error",
}

func Respond(c *gin.Context, code StatusCode, data interface{}) {
	c.JSON(200, gin.H{
		"code": code,
		"data": data,
	})
}

func RespondError(c *gin.Context, code StatusCode) {
	c.JSON(400, gin.H{
		"code": code,
		"msg":  ErrorMsgMap[code],
	})
}

func RespondSuccess(c *gin.Context, data interface{}) {
	Respond(c, OK, data)
}
