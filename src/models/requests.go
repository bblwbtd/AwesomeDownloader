package models

type DownloadRequest struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}
