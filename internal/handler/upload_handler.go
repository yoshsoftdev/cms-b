package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storageBaseURL string
}

func NewUploadHandler(storageBaseURL string) *UploadHandler {
	return &UploadHandler{storageBaseURL: storageBaseURL}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	ext := filepath.Ext(file.Filename)
	url := fmt.Sprintf("%s/%d%s", h.storageBaseURL, time.Now().UnixNano(), ext)
	c.JSON(http.StatusOK, gin.H{"url": url, "filename": file.Filename})
}
