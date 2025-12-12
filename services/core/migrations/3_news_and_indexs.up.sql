-- Add category column to news table
ALTER TABLE news ADD COLUMN category VARCHAR(50) DEFAULT 'general';

-- Add status column to channels table
ALTER TABLE channels ADD COLUMN status VARCHAR(15) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'banned'));
ALTER TABLE channels ADD COLUMN created_by INT DEFAULT NULL REFERENCES users(id) ON DELETE SET NULL;

-- Add index for category on news table
CREATE INDEX IF NOT EXISTS idx_news_category ON news(category);

-- Add index for created_at on news table
CREATE INDEX IF NOT EXISTS idx_news_created_at ON news(created_at DESC);

-- Add index for created_by on channels table
CREATE INDEX IF NOT EXISTS idx_channel_created_by ON channels(created_by);

-- Add index in status on channels table
CREATE INDEX IF NOT EXISTS idx_channel_status ON channels(status);