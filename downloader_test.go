package main

import (
    "AwesomeDownloader/entities"
    "os"
    "path"
    "testing"
)

func TestDownload(t *testing.T) {
    task := &entities.DownloadTask{
        URL:  "https://pic.netbian.com/uploads/allimg/170424/104135-14930016950de4.jpg",
        Path: path.Join("temp", "test.jpg"),
    }

    err := NewDownloader().Download(task, nil)
    if err != nil {
        t.Error(err)
        return
    }

    _, err = os.Stat(task.Path)
    if os.IsNotExist(err) {
        t.Error(err)
    }
}
