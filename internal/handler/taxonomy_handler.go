package handler

import (
	"net/http"
	"strconv"

	"cms/internal/model"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type TaxonomyHandler struct {
	service *service.TaxonomyService
}

func NewTaxonomyHandler(s *service.TaxonomyService) *TaxonomyHandler {
	return &TaxonomyHandler{service: s}
}

func (h *TaxonomyHandler) List(c *gin.Context) {
	typ := model.TaxonomyType(c.DefaultQuery("type", string(model.TaxonomyKeyword)))
	if typ != model.TaxonomyKeyword && typ != model.TaxonomyCategory {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid taxonomy type"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := h.service.List(typ, c.Query("search"), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}
