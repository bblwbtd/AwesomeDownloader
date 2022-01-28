package core

type DownloadStatus = string

const (
	Pending     DownloadStatus = "Pending"
	Downloading DownloadStatus = "Downloading"
	Paused      DownloadStatus = "Paused"
	Canceled    DownloadStatus = "Canceled"
	Finished    DownloadStatus = "Finished"
	Error       DownloadStatus = "Error"
	Unknown     DownloadStatus = "Unknown"
)
