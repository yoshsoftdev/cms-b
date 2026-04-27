package db

import (
	"fmt"
	"log"

	"cms/config"
	"cms/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.DBDriver {
	case "postgres":
		dialector = postgres.Open(cfg.DBDSN)
	default:
		dialector = sqlite.Open(cfg.DBDSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	if err := db.AutoMigrate(&model.Post{}, &model.Page{}, &model.Navigation{}, &model.TaxonomyTerm{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	if err := seedData(db, cfg); err != nil {
		return nil, fmt.Errorf("seed data: %w", err)
	}

	return db, nil
}

func seedData(db *gorm.DB, cfg config.Config) error {
	var postCount int64
	if err := db.Model(&model.Post{}).Count(&postCount).Error; err != nil {
		return err
	}
	if postCount == 0 {
		post := model.Post{
			Title:           "Welcome to CMS",
			Slug:            "welcome-to-cms",
			Content:         "<p>Hello from the headless CMS.</p>",
			Excerpt:         "Kickstart content.",
			Status:          model.StatusPublished,
			MetaTitle:       "Welcome",
			MetaDescription: "Welcome page",
		}
		if err := db.Create(&post).Error; err != nil {
			return err
		}
	}

	var pageCount int64
	if err := db.Model(&model.Page{}).Count(&pageCount).Error; err != nil {
		return err
	}
	if pageCount == 0 {
		page := model.Page{
			Title:           "About",
			Slug:            "about",
			Content:         "<p>About us page.</p>",
			Excerpt:         "About page",
			Status:          model.StatusPublished,
			MetaTitle:       "About",
			MetaDescription: "About page",
		}
		if err := db.Create(&page).Error; err != nil {
			return err
		}
	}

	var navCount int64
	if err := db.Model(&model.Navigation{}).Count(&navCount).Error; err != nil {
		return err
	}
	if navCount == 0 {
		navs := []model.Navigation{
			{Label: "Home", URL: "/", Order: 1, Type: "internal", Target: "_self", IsActive: true, Location: "header"},
			{Label: "About", URL: "/about", Order: 2, Type: "internal", Target: "_self", IsActive: true, Location: "header"},
			{Label: "Privacy", URL: "/privacy", Order: 1, Type: "internal", Target: "_self", IsActive: true, Location: "footer"},
		}
		if err := db.Create(&navs).Error; err != nil {
			return err
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	log.Printf("seed admin credentials: username=%s password=%s hash=%s", cfg.AdminUser, cfg.AdminPass, string(hash))
	return nil
}
