package service

import (
	"cms/internal/model"
	"cms/internal/repository"
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ContentService struct {
	repo             *repository.ContentRepository
	deployScriptPath string
}

func NewContentService(repo *repository.ContentRepository, deployScriptPath string) *ContentService {
	return &ContentService{repo: repo, deployScriptPath: strings.TrimSpace(deployScriptPath)}
}

func (s *ContentService) ListByLocale(q repository.ContentQuery) ([]repository.ContentListItem, int64, error) {
	return s.repo.ListByLocale(q)
}

func (s *ContentService) GetByID(id, locale string) (*repository.ContentDetail, error) {
	return s.repo.GetByIDWithLocale(id, locale)
}

func (s *ContentService) GetBySlug(slug, locale string) (*repository.ContentDetail, error) {
	return s.repo.GetBySlugWithLocale(slug, locale)
}

type CreateContentRequest struct {
	Slug         string                     `json:"slug"`
	Type         string                     `json:"type"`
	Status       string                     `json:"status"`
	AuthorID     string                     `json:"author_id"`
	PublishAt    *string                    `json:"publish_at"`
	Translations []CreateContentTranslation `json:"translations"`
	TaxonomyIDs  []string                   `json:"taxonomy_ids"`
	MediaLinks   []CreateContentMediaLink   `json:"media_links"`
}

type CreateContentTranslation struct {
	Locale          string `json:"locale"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	ContentJSON     string `json:"content_json"`
	Excerpt         string `json:"excerpt"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
	Keywords        string `json:"keywords"`
}

type CreateContentMediaLink struct {
	MediaID string `json:"media_id"`
	Type    string `json:"type"`
}

func (s *ContentService) Create(req CreateContentRequest) error {
	return s.persist("", req)
}

func (s *ContentService) Update(id string, req CreateContentRequest) error {
	return s.persist(id, req)
}

func (s *ContentService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ContentService) persist(id string, req CreateContentRequest) error {
	contentType := model.ContentType(strings.ToLower(strings.TrimSpace(req.Type)))
	if contentType != model.ContentTypePost && contentType != model.ContentTypePage {
		contentType = model.ContentTypePost
	}
	status := model.ContentStatus(strings.ToLower(strings.TrimSpace(req.Status)))
	if status != model.StatusDraft && status != model.StatusPublished {
		status = model.StatusDraft
	}

	content := model.ContentRecord{
		Slug:   strings.TrimSpace(req.Slug),
		Type:   contentType,
		Status: status,
	}

	if req.AuthorID != "" {
		if id, err := uuid.Parse(req.AuthorID); err == nil {
			content.AuthorID = &id
		}
	}
	if req.PublishAt != nil && *req.PublishAt != "" {
		if parsed, err := time.Parse(time.RFC3339, *req.PublishAt); err == nil {
			content.PublishAt = &parsed
		}
	}

	translations := make([]model.ContentTranslation, 0, len(req.Translations))
	for _, t := range req.Translations {
		locale := model.LocaleCode(strings.ToLower(strings.TrimSpace(t.Locale)))
		if locale != model.LocaleEN && locale != model.LocaleID {
			continue
		}
		translations = append(translations, model.ContentTranslation{
			Locale:          locale,
			Title:           strings.TrimSpace(t.Title),
			Content:         strings.TrimSpace(t.Content),
			ContentJSON:     strings.TrimSpace(t.ContentJSON),
			Excerpt:         strings.TrimSpace(t.Excerpt),
			MetaTitle:       strings.TrimSpace(t.MetaTitle),
			MetaDescription: strings.TrimSpace(t.MetaDescription),
			Keywords:        strings.TrimSpace(t.Keywords),
		})
	}

	taxonomyIDs := make([]uuid.UUID, 0, len(req.TaxonomyIDs))
	for _, raw := range req.TaxonomyIDs {
		if id, err := uuid.Parse(raw); err == nil {
			taxonomyIDs = append(taxonomyIDs, id)
		}
	}
	mediaLinks := make([]model.ContentMedia, 0, len(req.MediaLinks))
	for _, m := range req.MediaLinks {
		id, err := uuid.Parse(m.MediaID)
		if err != nil {
			continue
		}
		mediaLinks = append(mediaLinks, model.ContentMedia{MediaID: id, Type: strings.TrimSpace(m.Type)})
	}

	input := repository.CreateContentInput{
		Content:      content,
		Translations: translations,
		TaxonomyIDs:  taxonomyIDs,
		MediaLinks:   mediaLinks,
	}
	var err error
	if id == "" {
		err = s.repo.Create(input)
	} else {
		err = s.repo.Update(id, input)
	}
	if err != nil {
		return err
	}
	if content.Status == model.StatusPublished {
		s.runDeployScript()
	}
	return nil
}

func (s *ContentService) runDeployScript() {
	if s.deployScriptPath == "" {
		return
	}
	if filepath.Ext(s.deployScriptPath) != ".sh" {
		log.Printf("deploy skipped: DEPLOY_SCRIPT_PATH must point to .sh file, got: %s", s.deployScriptPath)
		return
	}
	go func(scriptPath string) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, scriptPath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("deploy script failed (%s): %v\n%s", scriptPath, err, string(out))
			return
		}
		log.Printf("deploy script executed: %s", scriptPath)
	}(s.deployScriptPath)
}
