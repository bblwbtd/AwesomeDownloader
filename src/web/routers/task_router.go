package routers

import (
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/handlers"
	"AwesomeDownloader/src/web/models"
	"github.com/gin-gonic/gin"
)

func mountTaskRouter(router *gin.RouterGroup) {
	task := router.Group("/task")

	task.POST("/add", addTask)
	task.POST("/remove/:id", removeTask)
	task.POST("/pause/:id", pauseTask)
	task.POST("/unpause/:id", unpauseTask)
	task.POST("/cancel/:id", cancelTask)
}

func addTask(ctx *gin.Context) {
	req := new(models.DownloadRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidBody, nil)
		return
	}

	task := handlers.AddTasks(downloader, req.Tasks)

	utils.RespondSuccess(ctx, task)
}

func removeTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID, err)
		return
	}

	err = handlers.CancelTasks(downloader, []uint{id})
	if err != nil {
		utils.RespondError(ctx, utils.RuntimeError, err)
		return
	}

	utils.RespondSuccess(ctx, "")
}

func pauseTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID, err)
		return
	}

	err = handlers.PauseTask(downloader, []uint{id})
	if err != nil {
		utils.RespondError(ctx, utils.RuntimeError, err)
		return
	}

	utils.RespondSuccess(ctx, "")
}

func unpauseTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID, err)
		return
	}

	err = handlers.UnpauseTask(downloader, []uint{id})
	if err != nil {
		utils.RespondError(ctx, utils.RuntimeError, err)
		return
	}

	utils.RespondSuccess(ctx, "")
}

func cancelTask(ctx *gin.Context) {
	id, err := utils.ExtractID(ctx)
	if err != nil {
		utils.RespondError(ctx, utils.InvalidID, err)
		return
	}

	err = handlers.CancelTasks(downloader, []uint{id})
	if err != nil {
		utils.RespondError(ctx, utils.RuntimeError, err)
		return
	}

	utils.RespondSuccess(ctx, "")
}
