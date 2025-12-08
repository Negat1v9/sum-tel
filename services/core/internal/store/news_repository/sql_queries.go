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

	countNewsByUserSourcesQuary = `
		SELECT count(DISTINCT ns.news_id) as total FROM user_subscriptions us
				JOIN news_sources ns ON ns.channel_id = us.channel_id
			WHERE us.user_id = $1
	`
	getNewsByUserSourcesQuary = `
		SELECT n.id, n.title, n.summary, n.language, n.created_at, COUNT(ns.id) AS number_of_sources
				FROM user_subscriptions us
			JOIN news_sources ns ON ns.channel_id = us.channel_id
			JOIN news n ON ns.news_id = n.id
				WHERE us.user_id = $1
			GROUP BY n.id
				ORDER BY n.created_at DESC
			LIMIT $2 OFFSET $3
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

	getNewsSourcesByNewsIDQuery = `
			SELECT ns.id, ns.message_id, ns.news_id, ns.channel_id, c.username AS channel_name
				FROM news_sources ns
			JOIN channels c ON c.id = ns.channel_id
			WHERE news_id = $1
	`
)
