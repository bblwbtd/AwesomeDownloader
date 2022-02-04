package models

type TaskMeta struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type DownloadRequest struct {
	Tasks []*TaskMeta `json:"tasks"`
}
