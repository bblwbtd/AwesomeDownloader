package handlers

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/downloader"
	"AwesomeDownloader/src/web/models"
	"context"
	"path"
	"sync"
)

var taskChannel = make(chan *entities.DownloadTask, 1024)
var cancellations = sync.Map{}
var downloadProgress = sync.Map{}
var dl = downloader.NewDownloader()

func StartScheduler() {
	cfg := config.GetConfig()

	for i := 0; i < cfg.MaxConnections; i += 1 {
		go func() {
			for {
				task := <-taskChannel
				database.DB.Take(task)
				if task.Status != entities.Pending {
					continue
				}
				ctx, cancel := context.WithCancel(context.TODO())
				cancellations.Store(task.ID, cancel)
				options := &downloader.DownloadOptions{
					UpdateSize: func(size uint64) {
						task.Size = size
						database.DB.Save(task)
					},
					OnProgress: func(size uint64) {
						go downloadProgress.Store(task.ID, size)
					},
				}

				task.Status = entities.Downloading
				database.DB.Save(task)
				if err := dl.Download(ctx, task, options); err != nil {
					task.Status = entities.Error
					database.DB.Save(task)
				}

				task.Status = entities.Finished
				database.DB.Save(task)

				cancellations.Delete(task.ID)
			}
		}()
	}

	var pendingTasks []entities.DownloadTask
	database.DB.Where("status = ?", entities.Pending).Find(&pendingTasks)
	for _, task := range pendingTasks {
		taskChannel <- &task
	}
}

func AddTask(request *models.DownloadRequest) *entities.DownloadTask {
	cfg := config.GetConfig()

	task := &entities.DownloadTask{
		URL:    request.URL,
		Path:   path.Join(cfg.DownloadDir, request.Path),
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
