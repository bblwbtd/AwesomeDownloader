package core

import (
	"AwesomeDownloader/src/database"
	"AwesomeDownloader/src/database/entities"
	"context"
	"log"
	"os"
)

type taskDecorator struct {
	entity         *entities.Task
	cancel         context.CancelFunc
	ctx            context.Context
	downloadedSize uint64
	retryCount     int

}

func NewDecoratedTask(entity *entities.Task) *taskDecorator {
	ctx, cancel := context.WithCancel(context.TODO())

	return &taskDecorator{entity: entity, cancel: cancel, ctx: ctx}
}

func (d *taskDecorator) GetDownloadedSize() uint64 {
	return d.downloadedSize
}

func (d *taskDecorator) SetDownloadedSize(downloadedSize uint64) {
	d.downloadedSize = downloadedSize
}

func (d *taskDecorator) GetTaskStatus() (DownloadStatus, error) {
	if tx := database.DB.Take(d.entity); tx.Error != nil {
		return Unknown, tx.Error
	}

	return d.entity.Status, nil
}

func (d *taskDecorator) SetTaskStatus(status DownloadStatus) error {
	d.entity.Status = status
	if tx := database.DB.Save(d.entity); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (d *taskDecorator) SetTaskSize(size uint64) error {
	d.entity.Size = size
	if tx := database.DB.Save(d.entity); tx.Error != nil {
		return tx.Error
	}

	return nil

}

func (d *taskDecorator) Cancel() {
	if d.cancel != nil {
		d.cancel()
	}

	err := os.Remove(d.entity.Path)
	if err != nil {
		log.Println("Error occur while removing:", d.entity.Path)
		return
	}
}
