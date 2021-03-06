package handlers

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/downloader"
	"AwesomeDownloader/src/utils"
	"AwesomeDownloader/src/web/models"
	"context"
)

func AddTask(request *models.DownloadRequest) *entities.DownloadTask {

	task := &entities.DownloadTask{
		URL:    request.URL,
		Path:   utils.GetDownloadPath(request.Path),
		Status: entities.Pending,
	}
	database.DB.Create(task)
	go func() {
		downloader.TaskChannel <- task
	}()

	return task
}

func RemoveTask(id uint) {
	downloader.DownloadProgress.Delete(id)
	cancel(id)
	database.DB.Delete(&entities.DownloadTask{}, id)
}

func PauseTask(id uint) {
	cancel(id)
	database.DB.Model(&entities.DownloadTask{}).Where("id = ?", id).Update("status", entities.Paused)
}

func UnPauseTask(id uint) {
	task := new(entities.DownloadTask)
	err := database.DB.Take(task, id).Error
	if err != nil {
		return
	}

	task.Status = entities.Pending
	database.DB.Save(task)

	go func() {
		downloader.TaskChannel <- task
	}()
}

func CancelTask(id uint) {
	cancel(id)

	database.DB.Model(&entities.DownloadTask{}).Where("id = ?", id).Update("status", entities.Canceled)
}

func cancel(id uint) {
	if cancel, loaded := downloader.Cancellations.LoadAndDelete(id); loaded {
		if cancelFun, ok := cancel.(context.CancelFunc); ok {
			cancelFun()
		}
	}
}
