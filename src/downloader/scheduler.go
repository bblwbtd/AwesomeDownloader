package downloader

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"Go/AwesomeDownloader/vendor/github.com/reactivex/rxgo/v2"
	"context"
	"github.com/reactivex/rxgo/v2"
	"log"
)

var (
	TaskChannel      = make(chan rxgo.Item)
	Cancellations    = make(map[uint]context.CancelFunc)
	DownloadProgress = make(map[uint]uint64)
	downloader       = NewDownloader()
)

func fetchTaskFromDB(_ context.Context, i interface{}) (interface{}, error) {

}

func observeTaskChannel() {
	cfg := config.GetConfig()

	rxgo.FromChannel(TaskChannel).Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		task := i.(*entities.Task)
		tx := database.DB.Take(task)
		if tx.Error != nil {
			task.Status = entities.Error
			log.Println(tx.Error.Error())
			return nil, tx.Error
		}
		return task, nil
	}).ForEach(func(i interface{}) {
		task := i.(*entities.Task)
		ctx, cancel := context.WithCancel(context.TODO())

	}, func(err error) {

	}, func() {

	}, rxgo.WithPool(cfg.MaxConnections))
}

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
				Cancellations[task.ID] = cancel
				options := &DownloadOptions{
					UpdateSize: func(size uint64) {
						task.Size = size
						database.DB.Save(task)
					},
					OnProgress: func(size uint64) {
						DownloadProgress[task.ID] = size
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

				delete(Cancellations, task.ID)
			}
		}()
	}

	resumeTasks()
}

func resumeTasks() {
	var pendingTasks []*entities.Task
	database.DB.Where("status = ?", entities.Pending).Find(&pendingTasks)
	for _, task := range pendingTasks {
		TaskChannel <- task
	}

	var downloadingTasks []*entities.Task
	database.DB.Where("status = ?", entities.Downloading).Find(&downloadingTasks)
	for _, task := range downloadingTasks {
		task.Status = entities.Pending
		database.DB.Save(task)
		TaskChannel <- task
	}
}
