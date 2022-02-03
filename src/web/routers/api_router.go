package routers

import "github.com/gin-gonic/gin"

func MountAPI(router *gin.Engine) {
	api := router.Group("/api")

	mountTaskRouter(api)
}
