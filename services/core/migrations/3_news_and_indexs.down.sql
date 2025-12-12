-- Remove index for category
DROP INDEX IF EXISTS idx_news_category;

-- Remove category column
ALTER TABLE news DROP COLUMN IF EXISTS category;