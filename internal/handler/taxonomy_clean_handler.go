package handler

import (
	"net/http"
	"strconv"

	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type TaxonomyCleanHandler struct {
	service *service.TaxonomyCleanService
}

func NewTaxonomyCleanHandler(s *service.TaxonomyCleanService) *TaxonomyCleanHandler {
	return &TaxonomyCleanHandler{service: s}
}

func (h *TaxonomyCleanHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := h.service.List(c.DefaultQuery("locale", "en"), c.Query("type"), c.Query("search"), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
