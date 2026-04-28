package main

import (
	"log"

	"cms/config"
	"cms/internal/handler"
	"cms/internal/middleware"
	"cms/internal/repository"
	"cms/internal/service"
	"cms/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	authService, err := service.NewAuthService(cfg.AdminUser, cfg.AdminPass, cfg.JWTSecret)
	if err != nil {
		log.Fatalf("auth init error: %v", err)
	}

	postService := service.NewPostService(repository.NewPostRepository(database))
	contentService := service.NewContentService(repository.NewContentRepository(database), cfg.DeployScriptPath)
	authorService := service.NewAuthorService(repository.NewAuthorRepository(database))
	mediaService := service.NewMediaService(repository.NewMediaRepository(database))
	taxonomyService := service.NewTaxonomyService(repository.NewTaxonomyRepository(database))
	taxonomyCleanService := service.NewTaxonomyCleanService(repository.NewTaxonomyCleanRepository(database))
	pageService := service.NewPageService(repository.NewPageRepository(database))
	navService := service.NewNavigationService(repository.NewGormNavigationRepository(database))

	authHandler := handler.NewAuthHandler(authService)
	postHandler := handler.NewPostHandler(postService, taxonomyService)
	contentHandler := handler.NewContentHandler(contentService)
	authorHandler := handler.NewAuthorHandler(authorService)
	mediaHandler := handler.NewMediaHandler(mediaService)
	taxonomyHandler := handler.NewTaxonomyHandler(taxonomyService)
	taxonomyCleanHandler := handler.NewTaxonomyCleanHandler(taxonomyCleanService)
	pageHandler := handler.NewPageHandler(pageService)
	navHandler := handler.NewNavigationHandler(navService)
	uploadHandler := handler.NewUploadHandler(cfg.StorageURL, cfg.StorageDomain, cfg.StorageDir)

	r := gin.Default()
	r.Static("/static", "./static")

	r.GET("/posts", postHandler.ListPublic)
	r.GET("/contents", contentHandler.List)
	r.GET("/contents/:id", contentHandler.GetByID)
	r.GET("/contents/slug/:slug", contentHandler.GetBySlug)
	r.GET("/taxonomies", taxonomyCleanHandler.List)
	r.GET("/authors", authorHandler.List)
	r.GET("/media", mediaHandler.List)
	r.GET("/posts/:slug", postHandler.GetBySlug)
	r.GET("/pages/:slug", pageHandler.GetBySlug)
	r.GET("/navigation", navHandler.GetTree)

	admin := r.Group("/admin")
	admin.POST("/auth/login", authHandler.Login)
	admin.Use(middleware.JWT(authService))

	admin.GET("/posts", postHandler.ListAdmin)
	admin.GET("/contents/:id", contentHandler.GetByID)
	admin.POST("/contents", contentHandler.Create)
	admin.PUT("/contents/:id", contentHandler.Update)
	admin.DELETE("/contents/:id", contentHandler.Delete)
	admin.GET("/posts/:id", postHandler.GetByID)
	admin.GET("/posts/:id/seo-audit", postHandler.GetSEOAuditByID)
	admin.GET("/posts/slug/:slug", postHandler.GetAdminBySlug)
	admin.GET("/posts/slug/:slug/seo-audit", postHandler.GetSEOAuditBySlug)
	admin.POST("/posts", postHandler.Create)
	admin.PUT("/posts/:id", postHandler.Update)
	admin.DELETE("/posts/:id", postHandler.Delete)
	admin.GET("/taxonomies", taxonomyHandler.List)
	admin.GET("/authors", authorHandler.List)
	admin.POST("/authors", authorHandler.Create)
	admin.PUT("/authors/:id", authorHandler.Update)
	admin.DELETE("/authors/:id", authorHandler.Delete)
	admin.GET("/media", mediaHandler.List)
	admin.POST("/media", mediaHandler.Create)
	admin.PUT("/media/:id", mediaHandler.Update)
	admin.DELETE("/media/:id", mediaHandler.Delete)

	admin.GET("/pages", pageHandler.ListAdmin)
	admin.GET("/pages/:id", pageHandler.GetByID)
	admin.POST("/pages", pageHandler.Create)
	admin.PUT("/pages/:id", pageHandler.Update)
	admin.DELETE("/pages/:id", pageHandler.Delete)

	admin.POST("/navigation", navHandler.Create)
	admin.PUT("/navigation/:id", navHandler.Update)
	admin.DELETE("/navigation/:id", navHandler.Delete)
	admin.POST("/navigation/reorder", navHandler.Reorder)

	admin.POST("/upload", uploadHandler.Upload)

	log.Printf("server running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
