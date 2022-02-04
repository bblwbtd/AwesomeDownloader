package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MergeHeaders(req *http.Request, header map[string]string) {
	for k, v := range header {
		req.Header.Add(k, v)
	}
}

func ExtractID(ctx *gin.Context) (uint, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	return uint(id), err
}
