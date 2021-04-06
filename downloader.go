package main

import (
    "AwesomeDownloader/entities"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strconv"
)

type DownloadOptions struct {
    updateSize func(size uint64)
    onProgress func(size uint64)
    header map[string]string
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
}

func NewDownloader() *Downloader {
    return &Downloader{}
}

var client = http.DefaultClient

func getContentLength(URL *url.URL) (uint64, error) {
    request, err := http.NewRequest("HEAD", URL.String(), nil)
    if err != nil {
        return 0, err
    }
    request.Header.Add("HOST", URL.Host)
    response, err := client.Do(request)
    if err != nil {
        return 0, err
    }
    contentLength, err := strconv.ParseUint(response.Header.Get("content-length"), 10, 64)
    if err != nil {
        return 0, err
    }
    return contentLength, nil
}

func mergeHeader(req *http.Request, header map[string]string)  {
    for k, v := range header {
        req.Header.Add(k, v)
    }
}

func (d *Downloader) Download(task *entities.DownloadTask, options *DownloadOptions) error {
    URL, err := url.Parse(task.URL)
    if err != nil {
        return err
    }

    length, err := getContentLength(URL)
    if err != nil {
        return err
    }
    if options != nil && options.updateSize != nil {
        options.updateSize(length)
    }

    downloadRequest, err := http.NewRequest("GET", task.URL, nil)
    if err != nil {
        return err
    }
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
        mergeHeader(downloadRequest, options.header)
    }
    response, err := client.Do(downloadRequest)
    if err != nil {
        return err
    }
    defer func() {
        _ = response.Body.Close()
    }()

    counter := WriteCounter{}
    if options != nil {
        counter.OnProgress = options.onProgress
    }

    _, err = io.Copy(file, io.TeeReader(response.Body, &counter))

    return err
}
