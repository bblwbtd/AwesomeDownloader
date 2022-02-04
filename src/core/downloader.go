package core

import (
	"AwesomeDownloader/src/config"
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type DownloadOptions struct {
	UpdateSize func(size uint64)
	OnProgress func(size uint64)
	header     map[string]string
}

type Downloader struct {
	client      *http.Client
	taskChannel chan rxgo.Item
	taskMap     map[uint]*TaskDecorator
}

func NewDownloader() *Downloader {
	downloader := &Downloader{
		client:      http.DefaultClient,
		taskChannel: make(chan rxgo.Item, 256),
		taskMap:     make(map[uint]*TaskDecorator),
	}

	downloader.SubscribeTaskChannel()

	return downloader
}

func (d *Downloader) SubscribeTaskChannel() {
	cfg := config.GetConfig()

	rxgo.FromChannel(d.taskChannel).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		decoratedTask := i.(*TaskDecorator)
		if err := decoratedTask.SetTaskStatus(Downloading); err != nil {
			return nil, err
		}

		d.taskMap[decoratedTask.entity.ID] = decoratedTask

		options := &DownloadOptions{
			UpdateSize: func(size uint64) {
				err := decoratedTask.SetTaskSize(size)
				if err != nil {
					log.Println("Fail to update task size:", err)
					return
				}
			},
			OnProgress: func(size uint64) {
				decoratedTask.SetDownloadedSize(size)
			},
		}

		if err := d.Download(decoratedTask.ctx, decoratedTask.entity, options); err != nil && err != context.Canceled {
			log.Println("Download error:", err)
			if err := decoratedTask.SetTaskStatus(Error); err != nil {
				log.Println("Failed to update task size:", err)
				return nil, err
			}
			return decoratedTask, nil
		}

		if err := decoratedTask.SetTaskStatus(Finished); err != nil {
			log.Println("Fail to update task status", err)
			return nil, err
		}

		return decoratedTask, nil
	}, rxgo.WithPool(cfg.MaxConnections)).ForEach(func(i interface{}) {
		decoratedTask := i.(*TaskDecorator)

		status, err := decoratedTask.GetTaskStatus()
		if err != nil {
			log.Println("Failed to access task status", err)
		}

		decoratedTask.cancel()

		if status == Error && decoratedTask.retryCount < cfg.MaxRetry {
			if err := decoratedTask.SetTaskStatus(Pending); err != nil {
				log.Println("Failed to update task status", err)
			}
			decoratedTask.retryCount += 1

		} else if status == Finished {
			if err := decoratedTask.SetTaskStatus(Finished); err != nil {
				log.Println("Failed to update task status", err)
			}
		}

	}, func(err error) {
		log.Println("Pipeline error:", err)
	}, func() {
		log.Println("Task channel closed")
	})
}

func (d *Downloader) CreateAndEnqueue(task *entities.Task) error {
	if err := database.DB.Create(task).Error; err != nil {
		return err
	}
	d.Enqueue(task)
	return nil
}

func (d *Downloader) Enqueue(task *entities.Task) {
	decoratedTask := NewDecoratedTask(task)
	d.taskChannel <- rxgo.Of(decoratedTask)
}

func (d *Downloader) getContentLength(URL *url.URL) (uint64, error) {
	request, err := http.NewRequest("HEAD", URL.String(), nil)
	if err != nil {
		return 0, err
	}
	request.Header.Add("HOST", URL.Host)
	response, err := d.client.Do(request)
	if err != nil {
		return 0, err
	}
	contentLength, err := strconv.ParseUint(response.Header.Get("content-length"), 10, 64)
	if err != nil {
		return 0, err
	}
	return contentLength, nil
}

func (d *Downloader) Download(ctx context.Context, task *entities.Task, options *DownloadOptions) error {
	log.Println("Begin downloading:", task.URL, task.Path)

	URL, err := url.Parse(task.URL)
	if err != nil {
		return err
	}

	length, err := d.getContentLength(URL)
	if err != nil {
		return err
	}
	if options != nil && options.UpdateSize != nil {
		options.UpdateSize(length)
	}

	downloadRequest, err := http.NewRequest("GET", task.URL, nil)
	if err != nil {
		return err
	}
	downloadRequest = downloadRequest.WithContext(ctx)
	var file *os.File
	stat, err := os.Stat(task.Path)
	if os.IsNotExist(err) {
		dir := filepath.Dir(task.Path)
		if err = os.MkdirAll(dir, 0777); err != nil {
			return err
		}
		file, err = os.OpenFile(task.Path, os.O_CREATE|os.O_WRONLY, 0777)
	} else {
		size := stat.Size()
		if uint64(size) == length {
			log.Println("File existed")
			return nil
		}
		file, err = os.OpenFile(task.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		downloadRequest.Header.Add("range", fmt.Sprintf("bytes=%d-%d", size, length))
	}
	if err != nil {
		return err
	}

	if options != nil && options.header != nil {
		utils.MergeHeader(downloadRequest, options.header)
	}
	response, err := d.client.Do(downloadRequest)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("body close error", err)
		}
	}(response.Body)

	stat, _ = file.Stat()
	counter := writeCounter{
		Size: uint64(stat.Size()),
	}
	if options != nil {
		counter.OnProgress = options.OnProgress
	}

	_, err = io.Copy(file, io.TeeReader(response.Body, &counter))

	if err != nil {
		return err
	}

	log.Println("Finish download:", task.URL, task.Path)
	return nil
}

func (d *Downloader) DeleteTasks(taskIDs []uint) error {
	tasks := make([]*entities.Task, len(taskIDs))

	for index, id := range taskIDs {
		task := d.taskMap[id]
		if task != nil {
			task.Cancel()
			delete(d.taskMap, id)
		}
		tasks[index] = task.entity
	}

	return database.DB.Delete(tasks).Error
}

func (d *Downloader) PauseTasks(taskIDs []uint) error {

	for _, id := range taskIDs {
		task := d.taskMap[id]
		if task != nil {
			task.Cancel()
		}
		log.Printf("Task %d is pasued\n", id)
	}

	return database.DB.Model(entities.Task{}).Where("id in ?", taskIDs).Update("status", Paused).Error
}

func (d *Downloader) UnpauseTasks(taskIDs []uint) error {
	var tasks []entities.Task

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&entities.Task{}).Where(taskIDs).Update("status", Pending)
		tx.Where(taskIDs).Find(&tasks)
		return nil
	}); err != nil {
		return err
	}

	for _, task := range tasks {
		d.Enqueue(&task)
	}

	return nil
}

func (d *Downloader) CancelTasks(taskIDs []uint) error {
	for _, id := range taskIDs {
		task := d.taskMap[id]
		if task != nil {
			task.Cancel()
		}
	}

	return database.DB.Model(&entities.Task{}).Where(taskIDs).Update("status", Canceled).Error
}
