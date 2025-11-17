package channel_repository

const (
	createChannelQuery = `
		INSERT INTO channels (id, username, title, description, parse_interval)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, title, description, parse_interval, last_parsed_at, created_at, updated_at
	`

	getChannelByIDQuery = `
		SELECT id, username, title, description, parse_interval, last_parsed_at, created_at, updated_at
		FROM channels
		WHERE id = $1
	`

	getChannelByUsernameQuery = `
		SELECT id, username, title, description, parse_interval, last_parsed_at, created_at, updated_at
		FROM channels
		WHERE username = $1
	`

	getAllChannelsQuery = `
		SELECT id, username, title, description, parse_interval, last_parsed_at, created_at, updated_at
		FROM channels
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	getUsernamesForParseQuery = `
		SELECT id, username, title, description, parse_interval, last_parsed_at, created_at, updated_at
		FROM channels
		WHERE last_parsed_at IS NULL OR (NOW() - INTERVAL '1 minute' * parse_interval > last_parsed_at)
		ORDER BY last_parsed_at ASC NULLS FIRST
		LIMIT $1 OFFSET $2
	`

	updateChannelQuery = `
		UPDATE channels
		SET username = $2, title = $3, description = $4, parse_interval = $5, last_parsed_at = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, title, description, parse_interval, last_parsed_at, created_at, updated_at
	`

	deleteChannelQuery = `
		DELETE FROM channels
		WHERE id = $1
		RETURNING id
	`
)
