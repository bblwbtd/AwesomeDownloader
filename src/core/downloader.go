package core

import (
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type DownloadOptions struct {
	UpdateSize func(task *entities.Task, size uint64)
	OnProgress func(task *entities.Task, size uint64)
	header     map[string]string
}

type Downloader struct {
	client      *http.Client
	taskChannel chan rxgo.Item
}

func NewDownloader() *Downloader {
	downloader := &Downloader{
		client: http.DefaultClient,
	}

	downloader.subscribeTaskChannel()

	return downloader
}

func (d *Downloader) subscribeTaskChannel() {
	rxgo.FromChannel(d.taskChannel).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		decoratedTask := i.(*TaskDecorator)
		if err := decoratedTask.setTaskStatus(Downloading); err != nil {
			return nil, err
		}

		context.WithCancel(context.TODO())

	}).ForEach(func(i interface{}) {

	}, func(err error) {
		log.Println(err)
	}, func() {
		log.Println("Task channel closed")
	})
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
	defer log.Println("Finish downloading:", task.URL, task.Path)
	URL, err := url.Parse(task.URL)
	if err != nil {
		return err
	}

	length, err := d.getContentLength(URL)
	if err != nil {
		return err
	}
	if options != nil && options.UpdateSize != nil {
		options.UpdateSize(task, length)
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
			log.Println("body close error")
		}
	}(response.Body)

	stat, _ = file.Stat()
	counter := writeCounter{
		Size: uint64(stat.Size()),
		task: task,
	}
	if options != nil {
		counter.OnProgress = options.OnProgress
	}

	_, err = io.Copy(file, io.TeeReader(response.Body, &counter))

	return err
}
