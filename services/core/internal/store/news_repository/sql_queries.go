package newsrepository

const (
	createNewsQuery = `
		INSERT INTO news (id, title, summary, language)
		VALUES ($1, $2, $3, $4)
	`

	getNewsByIDQuery = `
		SELECT id, title, summary, language, created_at
		FROM news
		WHERE id = $1
	`

	getAllNewsQuery = `
		SELECT id, title, summary, language, created_at
		FROM news
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	deleteNewsQuery = `
		DELETE FROM news
		WHERE id = $1
		RETURNING id, title, summary, language, created_at
	`

	createNewsSourceQuery = `
		INSERT INTO news_sources (news_id, message_id, channel_id)
		VALUES ($1, $2, $3)
	`

	createNewsSourcesQueryPrefix = `
		INSERT INTO news_sources (news_id, message_id, channel_id)
		VALUES %s
	`

	deleteNewsSourceQuery = `
		DELETE FROM news_sources
		WHERE id = $1
		RETURNING id, news_id, channel_id
	`

	deleteNewsSourcesByNewsIDQuery = `
		DELETE FROM news_sources
		WHERE news_id = $1
	`
)
