-- Migration: Add ASCII-normalized columns for search
ALTER TABLE posts ADD COLUMN title_ascii TEXT;
ALTER TABLE posts ADD COLUMN description_ascii TEXT;

-- Optional: Backfill existing data
UPDATE posts SET title_ascii = '';
UPDATE posts SET description_ascii = '';
