package handler

import (
	"net/http"

	"cms/internal/model"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type NavigationHandler struct {
	service *service.NavigationService
}

func NewNavigationHandler(s *service.NavigationService) *NavigationHandler {
	return &NavigationHandler{service: s}
}

func (h *NavigationHandler) GetTree(c *gin.Context) {
	location := c.Query("location")
	items, err := h.service.Tree(location, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *NavigationHandler) Create(c *gin.Context) {
	var item model.Navigation
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

func (h *NavigationHandler) Update(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "navigation not found"})
		return
	}
	if err := c.ShouldBindJSON(item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Update(item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *NavigationHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *NavigationHandler) Reorder(c *gin.Context) {
	var req struct {
		Location string   `json:"location" binding:"required"`
		IDs      []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Reorder(req.Location, req.IDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
