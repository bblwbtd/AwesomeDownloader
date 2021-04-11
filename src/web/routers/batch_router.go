package routers

import (
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/handlers"
	"AwesomeDownloader/src/web/models"
	"github.com/gin-gonic/gin"
)

func mountBatchRouter(router *gin.RouterGroup) {
	batchRouter := router.Group("/batch")

	batchRouter.POST("/add", addBatch)
	batchRouter.POST("/remove/:id", removeBatch)
	batchRouter.POST("/pause/:id", pauseBatch)
	batchRouter.POST("/unpause/:id", unpause)
}

func addBatch(ctx *gin.Context) {
	var batchRequest *models.BatchRequest
	err := ctx.BindJSON(batchRequest)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidBody)
		return
	}

	batch := handlers.AddBatch(batchRequest)

	utils.RespondSuccess(ctx, batch)
}

func removeBatch(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
	}

	handlers.RemoveBatch(id)

	utils.RespondSuccess(ctx, nil)
}

func pauseBatch(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
	}

	handlers.PauseBatch(id)

	utils.RespondSuccess(ctx, nil)
}

func unpause(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID)
	}

	handlers.UnPauseBatch(id)

	utils.RespondSuccess(ctx, nil)
}
