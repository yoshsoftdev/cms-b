PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contents (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL CHECK(type IN ('post','page')),
  slug TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft','published','archived')),
  author_id TEXT NULL,
  publish_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  FOREIGN KEY(author_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS content_translations (
  id TEXT PRIMARY KEY,
  content_id TEXT NOT NULL,
  locale TEXT NOT NULL CHECK(locale IN ('en','id')),
  title TEXT NOT NULL,
  content TEXT NULL,
  content_json TEXT NULL,
  excerpt TEXT NULL,
  meta_title TEXT NULL,
  meta_description TEXT NULL,
  keywords TEXT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_content_translations_content_locale ON content_translations(content_id, locale);

CREATE TABLE IF NOT EXISTS media (
  id TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  alt_text TEXT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS content_media (
  content_id TEXT NOT NULL,
  media_id TEXT NOT NULL,
  type TEXT NOT NULL,
  PRIMARY KEY (content_id, media_id, type),
  FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE,
  FOREIGN KEY(media_id) REFERENCES media(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS taxonomies (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL CHECK(type IN ('category','keyword','tag')),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS taxonomy_translations (
  id TEXT PRIMARY KEY,
  taxonomy_id TEXT NOT NULL,
  locale TEXT NOT NULL CHECK(locale IN ('en','id')),
  name TEXT NOT NULL,
  FOREIGN KEY(taxonomy_id) REFERENCES taxonomies(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_taxonomy_translations_taxonomy_locale ON taxonomy_translations(taxonomy_id, locale);

CREATE TABLE IF NOT EXISTS content_taxonomies (
  content_id TEXT NOT NULL,
  taxonomy_id TEXT NOT NULL,
  PRIMARY KEY (content_id, taxonomy_id),
  FOREIGN KEY(content_id) REFERENCES contents(id) ON DELETE CASCADE,
  FOREIGN KEY(taxonomy_id) REFERENCES taxonomies(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS navigations_v2 (
  id TEXT PRIMARY KEY,
  parent_id TEXT NULL,
  "order" INTEGER NOT NULL DEFAULT 0,
  type TEXT NOT NULL DEFAULT 'internal',
  target TEXT NOT NULL DEFAULT '_self',
  is_active INTEGER NOT NULL DEFAULT 1,
  location TEXT NOT NULL DEFAULT 'header',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(parent_id) REFERENCES navigations_v2(id)
);

CREATE TABLE IF NOT EXISTS navigation_translations (
  id TEXT PRIMARY KEY,
  navigation_id TEXT NOT NULL,
  locale TEXT NOT NULL CHECK(locale IN ('en','id')),
  label TEXT NOT NULL,
  url TEXT NOT NULL,
  FOREIGN KEY(navigation_id) REFERENCES navigations_v2(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_navigation_translations_nav_locale ON navigation_translations(navigation_id, locale);

CREATE INDEX IF NOT EXISTS idx_contents_type_status_publish_at ON contents(type, status, publish_at DESC);
CREATE INDEX IF NOT EXISTS idx_contents_author ON contents(author_id);
CREATE INDEX IF NOT EXISTS idx_contents_deleted_at ON contents(deleted_at);
CREATE INDEX IF NOT EXISTS idx_content_media_type ON content_media(type);
CREATE INDEX IF NOT EXISTS idx_taxonomies_type ON taxonomies(type);
CREATE INDEX IF NOT EXISTS idx_navigation_v2_parent_order ON navigations_v2(parent_id, "order");

COMMIT;
