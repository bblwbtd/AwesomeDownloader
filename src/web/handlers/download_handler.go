package handlers

import (
	config2 "AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/web/models"
	"context"
	"sync"
)

var taskChannel = make(chan *entities.DownloadTask, 1024)
var cancellations = sync.Map{}
var downloadProgress = sync.Map{}
var downloader = NewDownloader()

func StartScheduler() {
	config := config2.GetConfig()

	for i := 0; i < config.MaxConnections; i += 1 {
		go func() {
			for {
				task := <-taskChannel
				database.DB.Take(task)
				if task.Status != entities.Pending {
					continue
				}
				ctx, cancel := context.WithCancel(context.TODO())
				cancellations.Store(task.ID, cancel)
				options := &DownloadOptions{
					updateSize: func(size uint64) {
						task.Size = size
						database.DB.Save(task)
					},
					onProgress: func(size uint64) {
						go downloadProgress.Store(task.ID, size)
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

				cancellations.Delete(task.ID)
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

func RemoveTask(id uint) {
	downloadProgress.Delete(id)
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
		taskChannel <- task
	}()
}

func CancelTask(id uint) {
	cancel(id)

	database.DB.Model(&entities.DownloadTask{}).Where("id = ?", id).Update("status", entities.Canceled)
}

func cancel(id uint) {
	if cancel, loaded := cancellations.LoadAndDelete(id); loaded {
		if cancelFun, ok := cancel.(context.CancelFunc); ok {
			cancelFun()
		}
	}
}
