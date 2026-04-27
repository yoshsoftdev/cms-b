-- Base schema for PostgreSQL deployment
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(255) NOT NULL,
  title_id VARCHAR(255),
  slug VARCHAR(255) UNIQUE NOT NULL,
  content TEXT NOT NULL,
  content_id TEXT,
  content_json TEXT,
  content_json_id TEXT,
  excerpt TEXT,
  cover_image VARCHAR(500),
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  meta_title VARCHAR(255),
  meta_title_id VARCHAR(255),
  meta_description TEXT,
  meta_description_id TEXT,
  og_image VARCHAR(500),
  canonical_url VARCHAR(500),
  keywords TEXT,
  keywords_id TEXT,
  category VARCHAR(100),
  publish_date VARCHAR(40),
  author VARCHAR(255),
  read_time VARCHAR(40),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS pages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  content TEXT NOT NULL,
  excerpt TEXT,
  cover_image VARCHAR(500),
  status VARCHAR(20) NOT NULL DEFAULT 'draft',
  meta_title VARCHAR(255),
  meta_description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS taxonomy_terms (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  type VARCHAR(30) NOT NULL,
  name VARCHAR(120) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_taxonomy_type_name ON taxonomy_terms (type, name);

CREATE TABLE IF NOT EXISTS navigations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  label VARCHAR(255) NOT NULL,
  url VARCHAR(500) NOT NULL,
  parent_id UUID NULL,
  "order" INT NOT NULL DEFAULT 0,
  type VARCHAR(20) NOT NULL DEFAULT 'internal',
  target VARCHAR(20) NOT NULL DEFAULT '_self',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  location VARCHAR(20) NOT NULL DEFAULT 'header',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP NULL
);
