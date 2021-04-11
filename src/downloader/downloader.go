package downloader

import (
	"AwesomeDownloader/src/database/entities"
	"AwesomeDownloader/src/utils"
	"context"
	"fmt"
	"io"
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

type WriteCounter struct {
	Size       uint64
	OnProgress func(size uint64)
}

func (w *WriteCounter) Write(p []byte) (n int, err error) {
	length := len(p)
	w.Size = w.Size + uint64(length)
	if w.OnProgress != nil {
		w.OnProgress(uint64(length))
	}
	return length, nil
}

type Downloader struct {
	client *http.Client
}

func NewDownloader() *Downloader {
	downloader := &Downloader{
		client: http.DefaultClient,
	}
	return downloader
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

func (d *Downloader) Download(ctx context.Context, task *entities.DownloadTask, options *DownloadOptions) error {
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
	defer response.Body.Close()

	stat, _ = file.Stat()
	counter := WriteCounter{
		Size: uint64(stat.Size()),
	}
	if options != nil {
		counter.OnProgress = options.OnProgress
	}

	_, err = io.Copy(file, io.TeeReader(response.Body, &counter))

	return err
}
