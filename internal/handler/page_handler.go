package handler

import (
	"net/http"
	"strconv"

	"cms/internal/model"
	"cms/internal/repository"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service *service.PageService
}

func NewPageHandler(s *service.PageService) *PageHandler {
	return &PageHandler{service: s}
}

func (h *PageHandler) ListAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	items, total, err := h.service.List(repository.ContentFilter{Status: c.Query("status"), Search: c.Query("search"), Page: page, Limit: limit})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total, "page": page, "limit": limit})
}

func (h *PageHandler) GetBySlug(c *gin.Context) {
	item, err := h.service.GetBySlug(c.Param("slug"))
	if err != nil || item.Status != model.StatusPublished {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PageHandler) GetByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PageHandler) Create(c *gin.Context) {
	var item model.Page
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *PageHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}
	if err := c.ShouldBindJSON(existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(id, existing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

func (h *PageHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
