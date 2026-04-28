package handler

import (
	"fmt"
	"os"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storageBaseURL string
	storageDir     string
}

func NewUploadHandler(storageBaseURL, storageDir string) *UploadHandler {
	return &UploadHandler{
		storageBaseURL: storageBaseURL,
		storageDir:     storageDir,
	}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	if err := os.MkdirAll(h.storageDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare upload directory"})
		return
	}
	targetPath := filepath.Join(h.storageDir, filename)
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	url := fmt.Sprintf("%s/%s", h.storageBaseURL, filename)
	c.JSON(http.StatusOK, gin.H{"url": url, "filename": filename})
}
