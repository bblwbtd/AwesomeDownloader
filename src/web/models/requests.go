package models

type TaskMeta struct {
	URL     string            `json:"url"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

type DownloadRequest struct {
	Tasks []*TaskMeta `json:"tasks"`
}
