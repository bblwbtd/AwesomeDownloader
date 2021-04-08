package handlers

import (
	"AwesomeDownloader/src"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/models"
	"context"
)

var taskChannel = make(chan *entities.DownloadTask, 1024)
var cancellations = map[uint]context.CancelFunc{}
var downloadProgress = map[uint]uint64{}
var downloader = NewDownloader()

func StartScheduler() {
	config := src.GetConfig()

	for i := 0; i < config.MaxConnections; i += 1 {
		go func() {
			for {
				task := <-taskChannel
				ctx, cancel := context.WithCancel(context.TODO())
				cancellations[task.ID] = cancel
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

				delete(cancellations, task.ID)
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
	delete(downloadProgress, id)
	if cancel := cancellations[id]; cancel != nil {
		cancel()
		delete(cancellations, id)
	}
	database.DB.Delete(&entities.DownloadTask{}, id)
}

func PauseTask(id uint) {
	if cancel := cancellations[id]; cancel != nil {
		cancel()
		delete(cancellations, id)
	}

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

func CancelTask(id int) {
	if cancel := cancellations[uint(id)]; cancel != nil {
		cancel()
	}

	database.DB.Model(&entities.DownloadTask{}).Where("id = ?", id).Update("status", entities.Canceled)
}
