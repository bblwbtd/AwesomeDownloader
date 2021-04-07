package handlers

import (
	"AwesomeDownloader/src"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/entities"
	"AwesomeDownloader/src/models"
	"context"
)

var taskChannel = make(chan *entities.DownloadTask, 1024)
var contextMap = map[uint]context.CancelFunc{}
var downloadProgress = map[uint]uint64{}
var downloader = NewDownloader()

func StartScheduler() {
	config := src.GetConfig()

	for i := 0; i < config.MaxConnections; i += 1 {
		go func() {
			for {
				task := <-taskChannel
				ctx, cancel := context.WithCancel(context.TODO())
				contextMap[task.ID] = cancel
				options := &DownloadOptions{
					updateSize: func(size uint64) {
						task.Size = size
						database.DB.Save(task)
					},
					onProgress: func(size uint64) {
						downloadProgress[task.ID] = size
					},
				}

				task.Status = entities.Downloading
				database.DB.Save(task)
				if err := downloader.Download(ctx, task, options); err != nil {
					task.Status = entities.Error
					database.DB.Save(task)
				}

				task.Status = entities.Finished
				database.DB.Save(task)
			}
		}()
	}
}

func AddTask(request *models.DownloadRequest) *entities.DownloadTask {
	task := &entities.DownloadTask{
		URL:    request.URL,
		Path:   request.Path,
		Status: entities.Pending,
	}
	database.DB.Create(task)
	go func() {
		taskChannel <- task
	}()

	return task
}

func RemoveTask(id int) {
	delete(downloadProgress, uint(id))
	if cancel := contextMap[uint(id)]; cancel != nil {
		cancel()
	}
	database.DB.Delete(&entities.DownloadTask{}, id)
}

func PauseTask(id int) {
	if cancel := contextMap[uint(id)]; cancel != nil {
		cancel()
	}

	database.DB.Model(&entities.DownloadTask{}).Where("id = ?", id).Update("status", entities.Paused)
}

func unPauseTask() {

}

func CancelTask() {

}
