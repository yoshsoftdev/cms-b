INSERT INTO posts (id, title, slug, content, excerpt, status, meta_title, meta_description, created_at, updated_at)
VALUES
  (uuid_generate_v4(), 'Welcome to CMS', 'welcome-to-cms', '<p>Hello from the headless CMS.</p>', 'Kickstart content.', 'published', 'Welcome', 'Welcome page', NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;

INSERT INTO pages (id, title, slug, content, excerpt, status, meta_title, meta_description, created_at, updated_at)
VALUES
  (uuid_generate_v4(), 'About', 'about', '<p>About us page.</p>', 'About page', 'published', 'About', 'About page', NOW(), NOW())
ON CONFLICT (slug) DO NOTHING;
