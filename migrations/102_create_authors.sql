CREATE TABLE IF NOT EXISTS authors (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  photo TEXT DEFAULT '',
  alt_text TEXT DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO authors (id, name, photo, alt_text, created_at, updated_at)
SELECT
  u.id,
  u.name,
  '',
  '',
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
FROM users u
LEFT JOIN authors a ON lower(trim(a.name)) = lower(trim(u.name))
WHERE a.id IS NULL;
