package models

type DownloadRequest struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type BatchRequest struct {
	Name  string             `json:"name"`
	Tasks []*DownloadRequest `json:"tasks"`
}
