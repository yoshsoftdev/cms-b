package service

import (
	"fmt"
	"html"
	"strings"

	"cms/internal/model"
	"cms/internal/repository"
	"cms/internal/utils"
)

type PostService struct{ repo *repository.PostRepository }

func NewPostService(repo *repository.PostRepository) *PostService { return &PostService{repo: repo} }

func (s *PostService) Create(post *model.Post) error {
	post.Keywords = normalizeCSV(post.Keywords)
	post.KeywordsID = normalizeCSV(post.KeywordsID)
	post.Category = normalizeSingle(post.Category)
	post.Author = normalizeSingle(post.Author)
	post.ReadTime = normalizeSingle(post.ReadTime)
	post.PublishDate = normalizeSingle(post.PublishDate)
	post.TitleID = normalizeSingle(post.TitleID)
	post.OGImage = normalizeSingle(post.OGImage)
	post.CanonicalURL = normalizeSingle(post.CanonicalURL)
	post.ContentJSON = normalizeContent(post.ContentJSON)
	post.ContentJSONID = normalizeContent(post.ContentJSONID)
	post.Content = normalizeContent(post.Content)
	post.ContentID = normalizeContent(post.ContentID)
	post.MetaTitleID = normalizeSingle(post.MetaTitleID)
	post.MetaDescriptionID = normalizeSingle(post.MetaDescriptionID)
	base := utils.Slugify(post.Slug)
	if base == "" {
		base = utils.Slugify(post.Title)
	}
	for i := 0; i < 100; i++ {
		slug := base
		if i > 0 {
			slug = fmt.Sprintf("%s-%d", base, i)
		}
		exists, err := s.repo.ExistsBySlug(slug, "")
		if err != nil {
			return err
		}
		if !exists {
			post.Slug = slug
			return s.repo.Create(post)
		}
	}
	return fmt.Errorf("failed to generate unique slug")
}
func (s *PostService) Update(id string, post *model.Post) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	post.Keywords = normalizeCSV(post.Keywords)
	post.KeywordsID = normalizeCSV(post.KeywordsID)
	post.Category = normalizeSingle(post.Category)
	post.Author = normalizeSingle(post.Author)
	post.ReadTime = normalizeSingle(post.ReadTime)
	post.PublishDate = normalizeSingle(post.PublishDate)
	post.TitleID = normalizeSingle(post.TitleID)
	post.OGImage = normalizeSingle(post.OGImage)
	post.CanonicalURL = normalizeSingle(post.CanonicalURL)
	post.ContentJSON = normalizeContent(post.ContentJSON)
	post.ContentJSONID = normalizeContent(post.ContentJSONID)
	post.Content = normalizeContent(post.Content)
	post.ContentID = normalizeContent(post.ContentID)
	post.MetaTitleID = normalizeSingle(post.MetaTitleID)
	post.MetaDescriptionID = normalizeSingle(post.MetaDescriptionID)
	base := utils.Slugify(post.Slug)
	if base == "" {
		base = utils.Slugify(post.Title)
	}
	for i := 0; i < 100; i++ {
		slug := base
		if i > 0 {
			slug = fmt.Sprintf("%s-%d", base, i)
		}
		exists, err := s.repo.ExistsBySlug(slug, id)
		if err != nil {
			return err
		}
		if !exists {
			post.Slug = slug
			return s.repo.Update(post)
		}
	}
	return fmt.Errorf("failed to generate unique slug")
}
func (s *PostService) Delete(id string) error { return s.repo.Delete(id) }
func (s *PostService) GetByID(id string) (*model.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	post.Content = html.UnescapeString(post.Content)
	post.ContentJSON = html.UnescapeString(post.ContentJSON)
	post.ContentID = html.UnescapeString(post.ContentID)
	post.ContentJSONID = html.UnescapeString(post.ContentJSONID)
	return post, nil
}
func (s *PostService) GetBySlug(slug string) (*model.Post, error) {
	post, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	post.Content = html.UnescapeString(post.Content)
	post.ContentJSON = html.UnescapeString(post.ContentJSON)
	post.ContentID = html.UnescapeString(post.ContentID)
	post.ContentJSONID = html.UnescapeString(post.ContentJSONID)
	return post, nil
}
func (s *PostService) List(filter repository.ContentFilter) ([]model.Post, int64, error) {
	items, total, err := s.repo.List(filter)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Content = html.UnescapeString(items[i].Content)
		items[i].ContentJSON = html.UnescapeString(items[i].ContentJSON)
		items[i].ContentID = html.UnescapeString(items[i].ContentID)
		items[i].ContentJSONID = html.UnescapeString(items[i].ContentJSONID)
	}
	return items, total, nil
}

type SEOAudit struct {
	Warnings []string `json:"warnings"`
	Valid    bool     `json:"valid"`
}

func (s *PostService) BuildSEOAudit(post *model.Post, lang string) SEOAudit {
	metaTitle := post.MetaTitle
	metaDesc := post.MetaDescription
	keywords := post.Keywords
	title := post.Title
	content := post.Content
	if strings.ToLower(lang) == "id" {
		if post.MetaTitleID != "" {
			metaTitle = post.MetaTitleID
		}
		if post.MetaDescriptionID != "" {
			metaDesc = post.MetaDescriptionID
		}
		if post.KeywordsID != "" {
			keywords = post.KeywordsID
		}
		if post.TitleID != "" {
			title = post.TitleID
		}
		if post.ContentID != "" {
			content = post.ContentID
		}
	}
	warnings := make([]string, 0, 8)
	if strings.TrimSpace(post.Slug) == "" {
		warnings = append(warnings, "Slug is required")
	}
	if n := len(strings.TrimSpace(metaTitle)); n < 50 {
		warnings = append(warnings, "Meta title is too short (min 50)")
	} else if n > 60 {
		warnings = append(warnings, "Meta title is too long (max 60)")
	}
	if n := len(strings.TrimSpace(metaDesc)); n < 150 {
		warnings = append(warnings, "Meta description is too short (min 150)")
	} else if n > 160 {
		warnings = append(warnings, "Meta description is too long (max 160)")
	}
	if strings.TrimSpace(keywords) == "" {
		warnings = append(warnings, "Add at least one keyword")
	}
	if strings.TrimSpace(post.Category) == "" {
		warnings = append(warnings, "Add at least one category")
	}
	if strings.TrimSpace(post.Author) == "" {
		warnings = append(warnings, "Author is required")
	}
	if strings.TrimSpace(post.ReadTime) == "" {
		warnings = append(warnings, "Read time is required")
	}
	if strings.TrimSpace(title) == "" {
		warnings = append(warnings, "Title is required")
	}
	if strings.TrimSpace(content) == "" {
		warnings = append(warnings, "Content is required")
	}
	return SEOAudit{Warnings: warnings, Valid: len(warnings) == 0}
}

func normalizeSingle(v string) string {
	return html.EscapeString(v)
}

func normalizeContent(v string) string {
	return strings.TrimSpace(v)
}

func normalizeCSV(csv string) string {
	raw := strings.Split(csv, ",")
	seen := map[string]bool{}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		clean := strings.TrimSpace(item)
		if clean == "" {
			continue
		}
		key := strings.ToLower(clean)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, clean)
	}
	return html.EscapeString(strings.Join(out, ", "))
}
