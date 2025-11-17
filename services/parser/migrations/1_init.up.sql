CREATE TABLE raw_messages (
    id BIGSERIAL PRIMARY KEY,
    channel_id VARCHAR(255) NOT NULL,
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('text', 'image', 'text_image')),
    telegram_message_id BIGINT UNIQUE NOT NULL,
    html_text TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'new' CHECK (status IN ('new', 'processed')),
    media_urls TEXT[] DEFAULT '{}',
    message_date TIMESTAMP WITH TIME ZONE NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_raw_messages_channel_id ON raw_messages(channel_id);
CREATE INDEX idx_raw_messages_status ON raw_messages(status);
CREATE INDEX idx_raw_messages_telegram_channel_message_id ON raw_messages(channel_id, telegram_message_id);