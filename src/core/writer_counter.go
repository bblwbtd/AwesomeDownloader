package core

import "AwesomeDownloader/src/database/entities"

type writeCounter struct {
	Size       uint64
	task       *entities.Task
	OnProgress func(task *entities.Task, size uint64)
}

func (w *writeCounter) Write(p []byte) (n int, err error) {
	length := len(p)
	w.Size = w.Size + uint64(length)
	if w.OnProgress != nil {
		w.OnProgress(w.task, uint64(length))
	}
	return length, nil
}
