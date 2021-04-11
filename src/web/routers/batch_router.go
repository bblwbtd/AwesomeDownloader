package routers

import "github.com/gin-gonic/gin"

func mountBatchRouter(router *gin.Engine) {
	batchRouter := router.Group("/batch")

	batchRouter.GET("/j")
}
