package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storageBaseURL string
	storageDomain  string
	storageDir     string
}

func NewUploadHandler(storageBaseURL, storageDomain, storageDir string) *UploadHandler {
	return &UploadHandler{
		storageBaseURL: storageBaseURL,
		storageDomain:  strings.TrimSpace(storageDomain),
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
	url := fmt.Sprintf("%s/%s", h.resolveStorageBaseURL(c), filename)
	c.JSON(http.StatusOK, gin.H{"url": url, "filename": filename})
}

func (h *UploadHandler) resolveStorageBaseURL(c *gin.Context) string {
	if h.storageDomain != "" {
		return strings.TrimRight(h.storageDomain, "/") + "/static/uploads"
	}
	base := strings.TrimSpace(h.storageBaseURL)
	if base != "" && !strings.EqualFold(base, "auto") {
		return strings.TrimRight(base, "/")
	}
	scheme := c.Request.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	host := c.Request.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Request.Host
	}
	return fmt.Sprintf("%s://%s/static/uploads", scheme, host)
}
