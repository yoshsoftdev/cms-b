package handler

import (
	"net/http"
	"strconv"

	"cms/internal/repository"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	service *service.ContentService
}

func NewContentHandler(s *service.ContentService) *ContentHandler {
	return &ContentHandler{service: s}
}

func (h *ContentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	rows, total, err := h.service.ListByLocale(repository.ContentQuery{
		Locale: c.DefaultQuery("locale", "en"),
		Type:   c.Query("type"),
		Status: c.Query("status"),
		Search: c.Query("search"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "total": total, "page": page, "limit": limit})
}

func (h *ContentHandler) Create(c *gin.Context) {
	var req service.CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "content created"})
}

func (h *ContentHandler) GetByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Param("id"), c.DefaultQuery("locale", "en"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ContentHandler) GetBySlug(c *gin.Context) {
	item, err := h.service.GetBySlug(c.Param("slug"), c.DefaultQuery("locale", "en"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ContentHandler) Update(c *gin.Context) {
	var req service.CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(c.Param("id"), req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "content updated"})
}

func (h *ContentHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
