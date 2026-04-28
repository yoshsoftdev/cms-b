PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

DROP INDEX IF EXISTS idx_navigation_v2_parent_order;
DROP INDEX IF EXISTS idx_taxonomies_type;
DROP INDEX IF EXISTS idx_content_media_type;
DROP INDEX IF EXISTS idx_contents_deleted_at;
DROP INDEX IF EXISTS idx_contents_author;
DROP INDEX IF EXISTS idx_contents_type_status_publish_at;
DROP INDEX IF EXISTS idx_navigation_translations_nav_locale;
DROP INDEX IF EXISTS idx_taxonomy_translations_taxonomy_locale;
DROP INDEX IF EXISTS idx_content_translations_content_locale;

DROP TABLE IF EXISTS navigation_translations;
DROP TABLE IF EXISTS navigations_v2;
DROP TABLE IF EXISTS content_taxonomies;
DROP TABLE IF EXISTS taxonomy_translations;
DROP TABLE IF EXISTS taxonomies;
DROP TABLE IF EXISTS content_media;
DROP TABLE IF EXISTS media;
DROP TABLE IF EXISTS content_translations;
DROP TABLE IF EXISTS contents;
DROP TABLE IF EXISTS users;

COMMIT;
PRAGMA foreign_keys = ON;
