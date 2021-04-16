package downloader

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"context"
	"sync"
)

var (
	TaskChannel      = make(chan *entities.DownloadTask, 1024)
	Cancellations    = sync.Map{}
	DownloadProgress = sync.Map{}
	downloader       = NewDownloader()
)

func StartScheduler() {
	cfg := config.GetConfig()

	for i := 0; i < cfg.MaxConnections; i += 1 {
		go func() {
			for {
				task := <-TaskChannel
				database.DB.Take(task)
				if task.Status != entities.Pending {
					continue
				}
				ctx, cancel := context.WithCancel(context.TODO())
				Cancellations.Store(task.ID, cancel)
				options := &DownloadOptions{
					UpdateSize: func(size uint64) {
						task.Size = size
						database.DB.Save(task)
					},
					OnProgress: func(size uint64) {
						go DownloadProgress.Store(task.ID, size)
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

				Cancellations.Delete(task.ID)
			}
		}()
	}

	resumeTasks()
}

func resumeTasks() {
	var pendingTasks []*entities.DownloadTask
	database.DB.Where("status = ?", entities.Pending).Find(&pendingTasks)
	for _, task := range pendingTasks {
		TaskChannel <- task
	}

	var downloadingTasks []*entities.DownloadTask
	database.DB.Where("status = ?", entities.Downloading).Find(&downloadingTasks)
	for _, task := range downloadingTasks {
		task.Status = entities.Pending
		database.DB.Save(task)
		TaskChannel <- task
	}
}
