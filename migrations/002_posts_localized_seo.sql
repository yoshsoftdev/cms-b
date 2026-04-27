-- Adds bilingual/localized SEO fields required by v0 editor.
ALTER TABLE posts ADD COLUMN IF NOT EXISTS title_id VARCHAR(255);
ALTER TABLE posts ADD COLUMN IF NOT EXISTS content_id TEXT;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS content_json_id TEXT;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS meta_title_id VARCHAR(255);
ALTER TABLE posts ADD COLUMN IF NOT EXISTS meta_description_id TEXT;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS keywords_id TEXT;
