package core

type writeCounter struct {
	Size       uint64
	OnProgress func(size uint64)
}

func (w *writeCounter) Write(p []byte) (n int, err error) {
	length := len(p)
	w.Size = w.Size + uint64(length)
	if w.OnProgress != nil {
		w.OnProgress(uint64(length))
	}
	return length, nil
}
