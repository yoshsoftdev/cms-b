package main

import (
	"cms/config"
	"cms/internal/model"
	"cms/internal/utils"
	"cms/pkg/db"
	"errors"
	"flag"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type astroFrontmatter struct {
	Title           string `yaml:"title"`
	Slug            string `yaml:"slug"`
	Excerpt         string `yaml:"excerpt"`
	CoverImage      string `yaml:"cover_image"`
	Cover           string `yaml:"cover"`
	Status          string `yaml:"status"`
	Draft           *bool  `yaml:"draft"`
	MetaTitle       string `yaml:"meta_title"`
	MetaDescription string `yaml:"meta_description"`
	Description     string `yaml:"description"`
}

type astroContent struct {
	Frontmatter astroFrontmatter
	Body        string
	SourceFile  string
}

func main() {
	postsDir := flag.String("posts-dir", "", "Path to Astro posts directory (Markdown/MDX files)")
	pagesDir := flag.String("pages-dir", "", "Path to Astro pages directory (Markdown/MDX files)")
	dryRun := flag.Bool("dry-run", false, "Parse and show plan without writing to DB")
	flag.Parse()

	if *postsDir == "" && *pagesDir == "" {
		fmt.Println("At least one input is required: --posts-dir or --pages-dir")
		os.Exit(1)
	}

	cfg := config.Load()
	database, err := db.Connect(cfg)
	if err != nil {
		fmt.Printf("db error: %v\n", err)
		os.Exit(1)
	}

	if *postsDir != "" {
		posts, err := loadAstroFiles(*postsDir)
		if err != nil {
			fmt.Printf("load posts error: %v\n", err)
			os.Exit(1)
		}
		if err := migratePosts(database, posts, *dryRun); err != nil {
			fmt.Printf("migrate posts error: %v\n", err)
			os.Exit(1)
		}
	}

	if *pagesDir != "" {
		pages, err := loadAstroFiles(*pagesDir)
		if err != nil {
			fmt.Printf("load pages error: %v\n", err)
			os.Exit(1)
		}
		if err := migratePages(database, pages, *dryRun); err != nil {
			fmt.Printf("migrate pages error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Astro -> DB migration completed")
}

func migratePosts(database *gorm.DB, items []astroContent, dryRun bool) error {
	for _, item := range items {
		status := normalizeStatus(item.Frontmatter.Status, item.Frontmatter.Draft)
		title := fallback(item.Frontmatter.Title, deriveTitleFromFilename(item.SourceFile))
		slug := fallback(item.Frontmatter.Slug, utils.Slugify(title))
		excerpt := fallback(item.Frontmatter.Excerpt, item.Frontmatter.Description)
		cover := fallback(item.Frontmatter.CoverImage, item.Frontmatter.Cover)
		metaTitle := fallback(item.Frontmatter.MetaTitle, title)
		metaDescription := fallback(item.Frontmatter.MetaDescription, excerpt)
		content := html.EscapeString(item.Body)

		if dryRun {
			fmt.Printf("[DRY-RUN][POST] upsert slug=%s source=%s\n", slug, item.SourceFile)
			continue
		}

		var existing model.Post
		err := database.Where("slug = ?", slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			record := model.Post{
				Title:           title,
				Slug:            slug,
				Content:         content,
				Excerpt:         excerpt,
				CoverImage:      cover,
				Status:          model.ContentStatus(status),
				MetaTitle:       metaTitle,
				MetaDescription: metaDescription,
			}
			if err := database.Create(&record).Error; err != nil {
				return fmt.Errorf("create post %s: %w", slug, err)
			}
			fmt.Printf("[POST] created slug=%s\n", slug)
			continue
		}
		if err != nil {
			return fmt.Errorf("query post %s: %w", slug, err)
		}

		existing.Title = title
		existing.Content = content
		existing.Excerpt = excerpt
		existing.CoverImage = cover
		existing.Status = model.ContentStatus(status)
		existing.MetaTitle = metaTitle
		existing.MetaDescription = metaDescription
		if err := database.Save(&existing).Error; err != nil {
			return fmt.Errorf("update post %s: %w", slug, err)
		}
		fmt.Printf("[POST] updated slug=%s\n", slug)
	}
	return nil
}

func migratePages(database *gorm.DB, items []astroContent, dryRun bool) error {
	for _, item := range items {
		status := normalizeStatus(item.Frontmatter.Status, item.Frontmatter.Draft)
		title := fallback(item.Frontmatter.Title, deriveTitleFromFilename(item.SourceFile))
		slug := fallback(item.Frontmatter.Slug, utils.Slugify(title))
		excerpt := fallback(item.Frontmatter.Excerpt, item.Frontmatter.Description)
		cover := fallback(item.Frontmatter.CoverImage, item.Frontmatter.Cover)
		metaTitle := fallback(item.Frontmatter.MetaTitle, title)
		metaDescription := fallback(item.Frontmatter.MetaDescription, excerpt)
		content := html.EscapeString(item.Body)

		if dryRun {
			fmt.Printf("[DRY-RUN][PAGE] upsert slug=%s source=%s\n", slug, item.SourceFile)
			continue
		}

		var existing model.Page
		err := database.Where("slug = ?", slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			record := model.Page{
				Title:           title,
				Slug:            slug,
				Content:         content,
				Excerpt:         excerpt,
				CoverImage:      cover,
				Status:          model.ContentStatus(status),
				MetaTitle:       metaTitle,
				MetaDescription: metaDescription,
			}
			if err := database.Create(&record).Error; err != nil {
				return fmt.Errorf("create page %s: %w", slug, err)
			}
			fmt.Printf("[PAGE] created slug=%s\n", slug)
			continue
		}
		if err != nil {
			return fmt.Errorf("query page %s: %w", slug, err)
		}

		existing.Title = title
		existing.Content = content
		existing.Excerpt = excerpt
		existing.CoverImage = cover
		existing.Status = model.ContentStatus(status)
		existing.MetaTitle = metaTitle
		existing.MetaDescription = metaDescription
		if err := database.Save(&existing).Error; err != nil {
			return fmt.Errorf("update page %s: %w", slug, err)
		}
		fmt.Printf("[PAGE] updated slug=%s\n", slug)
	}
	return nil
}

func loadAstroFiles(baseDir string) ([]astroContent, error) {
	items := make([]astroContent, 0)
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".mdx" {
			return nil
		}
		item, err := parseAstroMarkdown(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		items = append(items, item)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func parseAstroMarkdown(path string) (astroContent, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return astroContent{}, err
	}
	content := string(b)

	if !strings.HasPrefix(content, "---\n") {
		return astroContent{Body: content, SourceFile: path}, nil
	}

	rest := strings.TrimPrefix(content, "---\n")
	parts := strings.SplitN(rest, "\n---\n", 2)
	if len(parts) != 2 {
		return astroContent{}, fmt.Errorf("invalid frontmatter format")
	}

	var fm astroFrontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), &fm); err != nil {
		return astroContent{}, err
	}

	return astroContent{Frontmatter: fm, Body: parts[1], SourceFile: path}, nil
}

func normalizeStatus(status string, draft *bool) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "draft" || normalized == "published" {
		return normalized
	}
	if draft != nil && *draft {
		return string(model.StatusDraft)
	}
	return string(model.StatusPublished)
}

func deriveTitleFromFilename(path string) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	return strings.TrimSpace(name)
}

func fallback(primary, secondary string) string {
	if strings.TrimSpace(primary) != "" {
		return strings.TrimSpace(primary)
	}
	return strings.TrimSpace(secondary)
}
