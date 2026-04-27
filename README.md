# Headless CMS Backend (Go + Gin)

Production-ready headless CMS backend for Astro SSG use-cases.

## Features
- Clean architecture: handler -> service -> repository
- JWT auth for admin API
- Posts and Pages with draft/published workflow
- Slug auto-generation + duplicate prevention
- Pagination + status filter + title search
- Navigation API with nested tree + reorder endpoint
- Media upload endpoint (returns external URL)
- SQLite fallback for development, PostgreSQL support
- GORM auto migration + sample seed data

## Project Structure
- `cmd/api` - app entrypoint
- `config` - environment config loader
- `internal/handler` - HTTP handlers
- `internal/service` - business logic
- `internal/repository` - data access
- `internal/model` - entities and DTO
- `internal/middleware` - JWT middleware
- `pkg/db` - db connection and migration setup
- `migrations` - SQL migration baseline
- `docs` - seed SQL and curl examples

## Run Locally
1. Copy env:
   - `cp .env.example .env`
2. Install dependencies:
   - `go mod tidy`
3. Run API:
   - `go run ./cmd/api`
4. Server runs at:
   - `http://localhost:8080`

## Astro to DB Migration
- Script migrasi ada di `cmd/migrate-astro/main.go`.
- Script membaca file `.md/.mdx` Astro (dengan/ tanpa frontmatter YAML), lalu upsert ke tabel `posts`/`pages` berdasarkan `slug`.

Contoh:
- Dry run:
  - `go run ./cmd/migrate-astro --posts-dir ../astro/src/content/blog --pages-dir ../astro/src/content/pages --dry-run`
- Eksekusi real:
  - `go run ./cmd/migrate-astro --posts-dir ../astro/src/content/blog --pages-dir ../astro/src/content/pages`

Mapping utama frontmatter:
- `title`, `slug`, `excerpt`, `cover_image`/`cover`, `status`, `draft`, `meta_title`, `meta_description`, `description`
- Jika `slug` kosong, akan dibuat otomatis dari `title`.
- Jika `status` kosong, default `published` (atau `draft` jika `draft: true`).

## Public Endpoints
- `GET /posts` (published only)
- `GET /posts/:slug`
- `GET /pages/:slug`
- `GET /navigation`

## Admin Endpoints
- `POST /admin/auth/login`
- `GET/POST/PUT/DELETE /admin/posts`
- `GET/POST/PUT/DELETE /admin/pages`
- `POST/PUT/DELETE /admin/navigation`
- `POST /admin/navigation/reorder`
- `POST /admin/upload`

## Notes for Astro Compatibility
- Slug-based routing ready (`/posts/:slug`, `/pages/:slug`)
- JSON responses are clean and SSG-friendly
- HTML content is escaped before persistence for safety

See `docs/curl-examples.md` for API usage samples.
