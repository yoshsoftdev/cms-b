package repository

import (
	"cms/internal/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContentRepository struct{ db *gorm.DB }

func NewContentRepository(db *gorm.DB) *ContentRepository { return &ContentRepository{db: db} }

func (r *ContentRepository) ResolveAuthorID(authorID, authorName string) (*uuid.UUID, error) {
	if parsed, err := uuid.Parse(strings.TrimSpace(authorID)); err == nil {
		return &parsed, nil
	}
	name := strings.TrimSpace(authorName)
	if name == "" {
		return nil, nil
	}
	email := strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@local.cms"
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err == nil {
		return &user.ID, nil
	}
	if err := r.db.Create(&model.User{Name: name, Email: email}).Error; err != nil {
		return nil, err
	}
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user.ID, nil
}

type ContentQuery struct {
	Locale string
	Type   string
	Status string
	Search string
	Page   int
	Limit  int
}

type ContentListItem struct {
	ID        uuid.UUID           `json:"id"`
	Type      model.ContentType   `json:"type"`
	Slug      string              `json:"slug"`
	Status    model.ContentStatus `json:"status"`
	PublishAt *time.Time          `json:"publish_at"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	AuthorID  *uuid.UUID          `json:"author_id"`
	Locale    model.LocaleCode    `json:"locale"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	ContentJSON string            `json:"content_json"`
	Excerpt   string              `json:"excerpt"`
	MetaTitle string              `json:"meta_title"`
	MetaDescription string        `json:"meta_description"`
	Keywords  string              `json:"keywords"`
}

func (r *ContentRepository) ListByLocale(q ContentQuery) ([]ContentListItem, int64, error) {
	if q.Locale == "" {
		q.Locale = string(model.LocaleEN)
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}

	base := r.db.Table("contents c").
		Joins("JOIN content_translations ct ON ct.content_id = c.id AND ct.locale = ?", q.Locale).
		Where("c.deleted_at IS NULL")

	if q.Type != "" {
		base = base.Where("c.type = ?", q.Type)
	}
	if q.Status != "" {
		base = base.Where("c.status = ?", q.Status)
	}
	if q.Search != "" {
		like := "%" + strings.TrimSpace(q.Search) + "%"
		base = base.Where("ct.title LIKE ? OR ct.meta_title LIKE ?", like, like)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []ContentListItem
	err := base.Select(
		"c.id, c.type, c.slug, c.status, c.publish_at, c.created_at, c.updated_at, c.author_id," +
			"ct.locale, ct.title, ct.content, ct.content_json, ct.excerpt, ct.meta_title, ct.meta_description, ct.keywords",
	).
		Order("c.created_at DESC").
		Offset((q.Page - 1) * q.Limit).
		Limit(q.Limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

type CreateContentInput struct {
	Content      model.ContentRecord
	Translations []model.ContentTranslation
	TaxonomyIDs  []uuid.UUID
	MediaLinks   []model.ContentMedia
}

func (r *ContentRepository) Create(input CreateContentInput) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&input.Content).Error; err != nil {
			return err
		}
		for i := range input.Translations {
			input.Translations[i].ContentID = input.Content.ID
		}
		if len(input.Translations) > 0 {
			if err := tx.Create(&input.Translations).Error; err != nil {
				return err
			}
		}
		for _, tID := range input.TaxonomyIDs {
			row := model.ContentTaxonomy{ContentID: input.Content.ID, TaxonomyID: tID}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		for _, link := range input.MediaLinks {
			link.ContentID = input.Content.ID
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type ContentAuthor struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Photo   string    `json:"photo"`
	AltText string    `json:"alt_text"`
}

type ContentTaxonomyItem struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	Name string    `json:"name"`
}

type ContentMediaItem struct {
	ID      uuid.UUID `json:"id"`
	URL     string    `json:"url"`
	AltText string    `json:"alt_text"`
	Type    string    `json:"type"`
}

type ContentDetail struct {
	ID        uuid.UUID           `json:"id"`
	Type      model.ContentType   `json:"type"`
	Slug      string              `json:"slug"`
	Status    model.ContentStatus `json:"status"`
	PublishAt *time.Time          `json:"publish_at"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	AuthorID  *uuid.UUID          `json:"author_id"`
	Author    *ContentAuthor      `json:"author,omitempty"`

	Translation  *model.ContentTranslation `json:"translation,omitempty"`
	Translations []model.ContentTranslation `json:"translations,omitempty"`
	Taxonomies   []ContentTaxonomyItem      `json:"taxonomies"`
	Media        []ContentMediaItem         `json:"media"`
}

func (r *ContentRepository) GetByIDWithLocale(id, locale string) (*ContentDetail, error) {
	var c model.ContentRecord
	if err := r.db.First(&c, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	out := &ContentDetail{
		ID:        c.ID,
		Type:      c.Type,
		Slug:      c.Slug,
		Status:    c.Status,
		PublishAt: c.PublishAt,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		AuthorID:  c.AuthorID,
		Taxonomies: []ContentTaxonomyItem{},
		Media:      []ContentMediaItem{},
	}

	if c.AuthorID != nil {
		var a model.Author
		if err := r.db.First(&a, "id = ?", *c.AuthorID).Error; err == nil {
			out.Author = &ContentAuthor{ID: a.ID, Name: a.Name, Photo: a.Photo, AltText: a.AltText}
		}
	}
	if locale == "" {
		locale = string(model.LocaleEN)
	}
	var tr model.ContentTranslation
	if err := r.db.First(&tr, "content_id = ? AND locale = ?", c.ID, locale).Error; err != nil {
		if locale != string(model.LocaleEN) {
			if errEn := r.db.First(&tr, "content_id = ? AND locale = ?", c.ID, string(model.LocaleEN)).Error; errEn == nil {
				out.Translation = &tr
			}
		}
	} else {
		out.Translation = &tr
	}
	_ = r.db.Where("content_id = ?", c.ID).Order("locale ASC").Find(&out.Translations).Error

	_ = r.db.Table("content_taxonomies ct").
		Select("t.id, t.type, tt.name").
		Joins("JOIN taxonomies t ON t.id = ct.taxonomy_id").
		Joins("JOIN taxonomy_translations tt ON tt.taxonomy_id = t.id AND tt.locale = ?", locale).
		Where("ct.content_id = ?", c.ID).
		Scan(&out.Taxonomies).Error

	_ = r.db.Table("content_media cm").
		Select("m.id, m.url, m.alt_text, cm.type").
		Joins("JOIN media m ON m.id = cm.media_id").
		Where("cm.content_id = ?", c.ID).
		Scan(&out.Media).Error

	return out, nil
}

func (r *ContentRepository) GetBySlugWithLocale(slug, locale string) (*ContentDetail, error) {
	var c model.ContentRecord
	if err := r.db.First(&c, "slug = ? AND deleted_at IS NULL", slug).Error; err != nil {
		return nil, err
	}
	return r.GetByIDWithLocale(c.ID.String(), locale)
}

func (r *ContentRepository) Delete(id string) error {
	return r.db.Delete(&model.ContentRecord{}, "id = ?", id).Error
}

func (r *ContentRepository) Update(id string, input CreateContentInput) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing model.ContentRecord
		if err := tx.First(&existing, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
			return err
		}
		input.Content.ID = existing.ID
		if err := tx.Model(&existing).Updates(map[string]any{
			"type":       input.Content.Type,
			"slug":       input.Content.Slug,
			"status":     input.Content.Status,
			"author_id":  input.Content.AuthorID,
			"publish_at": input.Content.PublishAt,
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("content_id = ?", existing.ID).Delete(&model.ContentTranslation{}).Error; err != nil {
			return err
		}
		for i := range input.Translations {
			input.Translations[i].ContentID = existing.ID
		}
		if len(input.Translations) > 0 {
			if err := tx.Create(&input.Translations).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("content_id = ?", existing.ID).Delete(&model.ContentTaxonomy{}).Error; err != nil {
			return err
		}
		for _, tID := range input.TaxonomyIDs {
			row := model.ContentTaxonomy{ContentID: existing.ID, TaxonomyID: tID}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("content_id = ?", existing.ID).Delete(&model.ContentMedia{}).Error; err != nil {
			return err
		}
		for _, link := range input.MediaLinks {
			link.ContentID = existing.ID
			if err := tx.Create(&link).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
