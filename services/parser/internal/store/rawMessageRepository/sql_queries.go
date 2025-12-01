package rawMessageRepository

const (
	createMessagesQuery = `
		INSERT INTO raw_messages (channel_id, content_type, telegram_message_id, html_text, status, media_urls, message_date, received_at)
		VALUES %s
		ON CONFLICT (channel_id, telegram_message_id) DO NOTHING
	`

	getChannelMessagesQuery = `
		SELECT id, channel_id, content_type, telegram_message_id, html_text, status, media_urls, message_date, received_at
		FROM raw_messages
		WHERE channel_id = $1
		ORDER BY message_date DESC
		LIMIT $2 OFFSET $3
	`

	getLatestChannelMessageQuery = `
		SELECT id, channel_id, content_type, telegram_message_id, html_text, status, media_urls, message_date, received_at
		FROM raw_messages
		WHERE channel_id = $1
		ORDER BY telegram_message_id DESC
		LIMIT 1
	`

	getAndProcessMessagesQuery = `
		WITH updated AS (
			UPDATE raw_messages
			SET status = 'processed'
			WHERE id IN (
				SELECT id FROM raw_messages
				WHERE status = 'new'
				ORDER BY telegram_message_id ASC
				LIMIT $1
			)
			RETURNING id, channel_id, content_type, telegram_message_id, html_text, status, media_urls, message_date, received_at
		)
		SELECT * FROM updated ORDER BY telegram_message_id ASC
	`

	updateMessagesStatusQuery = `
		UPDATE raw_messages
		SET status = 'processed'
		WHERE id = ANY($1)
	`
)
