CREATE TABLE IF NOT EXISTS news (
    id uuid PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    language VARCHAR(15) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_news_language ON news(language);

CREATE TABLE IF NOT EXISTS news_sources (
    id SERIAL PRIMARY KEY,
    message_id BIGINT NOT NULL,
    news_id uuid REFERENCES news(id) ON DELETE CASCADE,
    channel_id uuid REFERENCES channels(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_news_sources_news_id ON news_sources(news_id);
CREATE INDEX IF NOT EXISTS idx_news_sources_channel_id ON news_sources(channel_id);
