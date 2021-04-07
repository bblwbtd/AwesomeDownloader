package utils

import "net/http"

func MergeHeader(req *http.Request, header map[string]string) {
	for k, v := range header {
		req.Header.Add(k, v)
	}
}
