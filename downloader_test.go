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

    downloader := NewDownloader()
    err := downloader.Download(task, nil)
    if err != nil {
        t.Error(err)
        return
    }

    _, err = os.Stat(task.Path)
    if os.IsNotExist(err) {
        t.Error(err)
    }

}

func TestDownloadBreakPoint(t *testing.T) {
    task := &entities.DownloadTask{
        URL:  "https://pic.netbian.com/uploads/allimg/170424/104135-14930016950de4.jpg",
        Path: path.Join("temp", "test1.jpg"),
    }

    options := &DownloadOptions{
        header: map[string]string{
            "range": "bytes=0-150",
        },
    }

	downloader := NewDownloader()
	err := downloader.Download(task, options)
	if err != nil {
		t.Error(err)
		return
	}

	stat, err := os.Stat(task.Path)
	if err != nil {
		t.Error(err)
		return
	}

	if stat.Size() != 151 {
		t.Error("size is not correct")
		return 
	}

	err = downloader.Download(task, nil)
	if err != nil {
		t.Error(err)
		return
	}
	
}

func TestMain(m *testing.M) {
	os.RemoveAll("temp")
	code := m.Run()
	os.RemoveAll("temp")
	os.Exit(code)
}