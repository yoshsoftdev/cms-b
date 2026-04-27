package handler

import (
	"net/http"
	"strconv"

	"cms/internal/model"
	"cms/internal/repository"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service         *service.PostService
	taxonomyService *service.TaxonomyService
}

func NewPostHandler(s *service.PostService, ts *service.TaxonomyService) *PostHandler {
	return &PostHandler{service: s, taxonomyService: ts}
}

func (h *PostHandler) ListPublic(c *gin.Context) {
	h.list(c, string(model.StatusPublished))
}

func (h *PostHandler) ListAdmin(c *gin.Context) {
	h.list(c, c.Query("status"))
}

func (h *PostHandler) list(c *gin.Context, defaultStatus string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	if status == "" {
		status = defaultStatus
	}
	items, total, err := h.service.List(repository.ContentFilter{
		Status: status,
		Search: c.Query("search"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total, "page": page, "limit": limit})
}

func (h *PostHandler) GetBySlug(c *gin.Context) {
	item, err := h.service.GetBySlug(c.Param("slug"))
	if err != nil || item.Status != model.StatusPublished {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PostHandler) GetSEOAuditByID(c *gin.Context) {
	item, err := h.service.GetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	lang := c.DefaultQuery("lang", "en")
	c.JSON(http.StatusOK, h.service.BuildSEOAudit(item, lang))
}

func (h *PostHandler) GetAdminBySlug(c *gin.Context) {
	item, err := h.service.GetBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *PostHandler) GetSEOAuditBySlug(c *gin.Context) {
	item, err := h.service.GetBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	lang := c.DefaultQuery("lang", "en")
	c.JSON(http.StatusOK, h.service.BuildSEOAudit(item, lang))
}

func (h *PostHandler) Create(c *gin.Context) {
	var item model.Post
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if h.taxonomyService != nil {
		if err := h.taxonomyService.UpsertKeywordsFromCSV(item.Keywords); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.taxonomyService.UpsertCategory(item.Category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusCreated, item)
}

func (h *PostHandler) Update(c *gin.Context) {
	id := c.Param("id")
	existing, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
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
	if h.taxonomyService != nil {
		if err := h.taxonomyService.UpsertKeywordsFromCSV(existing.Keywords); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.taxonomyService.UpsertCategory(existing.Category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, existing)
}

func (h *PostHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
