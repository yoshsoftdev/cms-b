package repository

import (
	"sort"
	"strings"

	"gorm.io/gorm"
)

type TaxonomyCleanRepository struct{ db *gorm.DB }

func NewTaxonomyCleanRepository(db *gorm.DB) *TaxonomyCleanRepository {
	return &TaxonomyCleanRepository{db: db}
}

type TaxonomyCleanItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func taxonomyItemKey(t TaxonomyCleanItem) string {
	return strings.ToLower(strings.TrimSpace(t.Type)) + "\x00" + strings.ToLower(strings.TrimSpace(t.Name))
}

// mergeTaxonomyCleanAndLegacy combines normalized taxonomies with legacy taxonomy_terms rows (same CMS UI still writes to taxonomy_terms in some flows).
func mergeTaxonomyCleanAndLegacy(clean, legacy []TaxonomyCleanItem) []TaxonomyCleanItem {
	seen := make(map[string]struct{}, len(clean)+len(legacy))
	out := make([]TaxonomyCleanItem, 0, len(clean)+len(legacy))
	appendUnique := func(items []TaxonomyCleanItem) {
		for _, item := range items {
			item.Type = strings.TrimSpace(item.Type)
			item.Name = strings.TrimSpace(item.Name)
			if item.Name == "" {
				continue
			}
			k := taxonomyItemKey(item)
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			out = append(out, item)
		}
	}
	appendUnique(clean)
	appendUnique(legacy)
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

func (r *TaxonomyCleanRepository) List(locale, termType, search string, limit int) ([]TaxonomyCleanItem, error) {
	if locale == "" {
		locale = "en"
	}
	if limit <= 0 {
		limit = 20
	}
	searchTrim := strings.TrimSpace(search)
	like := "%" + searchTrim + "%"

	var clean []TaxonomyCleanItem
	q := r.db.Table("taxonomies t").
		Select("t.id, t.type, tt.name").
		Joins("JOIN taxonomy_translations tt ON tt.taxonomy_id = t.id AND tt.locale = ?", locale)
	if termType != "" {
		q = q.Where("t.type = ?", termType)
	}
	if searchTrim != "" {
		q = q.Where("tt.name LIKE ?", like)
	}
	if err := q.Order("tt.name ASC").Scan(&clean).Error; err != nil {
		return nil, err
	}

	var legacy []TaxonomyCleanItem
	lq := r.db.Table("taxonomy_terms").
		Select("id, type, name").
		Where("deleted_at IS NULL")
	if termType != "" {
		lq = lq.Where("type = ?", termType)
	}
	if searchTrim != "" {
		lq = lq.Where("name LIKE ?", like)
	}
	if err := lq.Order("name ASC").Scan(&legacy).Error; err != nil {
		return nil, err
	}

	merged := mergeTaxonomyCleanAndLegacy(clean, legacy)
	if len(merged) > limit {
		merged = merged[:limit]
	}
	return merged, nil
}
