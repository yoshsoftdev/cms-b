PRAGMA foreign_keys = ON;
BEGIN TRANSACTION;

INSERT INTO users (id, name, email, created_at, updated_at)
SELECT
  lower(hex(randomblob(4)))||'-'||lower(hex(randomblob(2)))||'-4'||substr(lower(hex(randomblob(2))),2)||'-a'||substr(lower(hex(randomblob(2))),2)||'-'||lower(hex(randomblob(6))),
  a.author_name,
  lower(replace(trim(a.author_name), ' ', '.')) || '@local.cms',
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM (
  SELECT DISTINCT TRIM(author) AS author_name FROM posts WHERE TRIM(COALESCE(author,'')) <> ''
) a
LEFT JOIN users u ON u.email = lower(replace(trim(a.author_name), ' ', '.')) || '@local.cms'
WHERE u.id IS NULL;

INSERT INTO users (id, name, email, created_at, updated_at)
SELECT '00000000-0000-4000-a000-000000000001', 'System', 'system@local.cms', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = '00000000-0000-4000-a000-000000000001');

INSERT INTO contents (id, type, slug, status, author_id, publish_at, created_at, updated_at, deleted_at)
SELECT
  p.id,
  'post',
  p.slug,
  CASE WHEN p.status IN ('draft','published') THEN p.status ELSE 'draft' END,
  COALESCE(
    (SELECT u.id FROM users u WHERE u.email = lower(replace(trim(COALESCE(p.author,'')), ' ', '.')) || '@local.cms' LIMIT 1),
    '00000000-0000-4000-a000-000000000001'
  ),
  CASE WHEN TRIM(COALESCE(p.publish_date,'')) = '' THEN NULL ELSE datetime(p.publish_date) END,
  p.created_at,
  p.updated_at,
  p.deleted_at
FROM posts p
LEFT JOIN contents c ON c.id = p.id
WHERE c.id IS NULL;

CREATE TABLE IF NOT EXISTS _legacy_page_id_map (
  old_page_id TEXT PRIMARY KEY,
  new_content_id TEXT NOT NULL UNIQUE
);

INSERT INTO _legacy_page_id_map (old_page_id, new_content_id)
SELECT
  p.id,
  lower(hex(randomblob(4)))||'-'||lower(hex(randomblob(2)))||'-4'||substr(lower(hex(randomblob(2))),2)||'-a'||substr(lower(hex(randomblob(2))),2)||'-'||lower(hex(randomblob(6)))
FROM pages p
LEFT JOIN _legacy_page_id_map m ON m.old_page_id = p.id
WHERE m.old_page_id IS NULL;

INSERT INTO contents (id, type, slug, status, author_id, publish_at, created_at, updated_at, deleted_at)
SELECT
  m.new_content_id,
  'page',
  pg.slug,
  CASE WHEN pg.status IN ('draft','published') THEN pg.status ELSE 'draft' END,
  '00000000-0000-4000-a000-000000000001',
  NULL,
  pg.created_at,
  pg.updated_at,
  pg.deleted_at
FROM pages pg
JOIN _legacy_page_id_map m ON m.old_page_id = pg.id
LEFT JOIN contents c ON c.id = m.new_content_id
WHERE c.id IS NULL;

INSERT INTO content_translations (id, content_id, locale, title, content, content_json, excerpt, meta_title, meta_description, keywords, created_at, updated_at)
SELECT
  lower(hex(randomblob(4)))||'-'||lower(hex(randomblob(2)))||'-4'||substr(lower(hex(randomblob(2))),2)||'-a'||substr(lower(hex(randomblob(2))),2)||'-'||lower(hex(randomblob(6))),
  p.id, 'en', p.title, p.content, p.content_json, p.excerpt, p.meta_title, p.meta_description, p.keywords, p.created_at, p.updated_at
FROM posts p
LEFT JOIN content_translations ct ON ct.content_id = p.id AND ct.locale = 'en'
WHERE ct.id IS NULL;

INSERT INTO content_translations (id, content_id, locale, title, content, content_json, excerpt, meta_title, meta_description, keywords, created_at, updated_at)
SELECT
  lower(hex(randomblob(4)))||'-'||lower(hex(randomblob(2)))||'-4'||substr(lower(hex(randomblob(2))),2)||'-a'||substr(lower(hex(randomblob(2))),2)||'-'||lower(hex(randomblob(6))),
  p.id, 'id',
  COALESCE(NULLIF(TRIM(p.title_id),''), p.title),
  COALESCE(NULLIF(TRIM(p.content_id),''), p.content),
  COALESCE(NULLIF(TRIM(p.content_json_id),''), p.content_json),
  p.excerpt,
  COALESCE(NULLIF(TRIM(p.meta_title_id),''), p.meta_title),
  COALESCE(NULLIF(TRIM(p.meta_description_id),''), p.meta_description),
  COALESCE(NULLIF(TRIM(p.keywords_id),''), p.keywords),
  p.created_at, p.updated_at
FROM posts p
LEFT JOIN content_translations ct ON ct.content_id = p.id AND ct.locale = 'id'
WHERE ct.id IS NULL;

COMMIT;
