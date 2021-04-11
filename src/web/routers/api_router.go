package routers

import "github.com/gin-gonic/gin"

func MountAPI(route *gin.Engine) {
	api := route.Group("/api")

	mountDownload(api)
}
