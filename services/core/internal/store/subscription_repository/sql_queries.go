package subscription_repository

const (
	createSubscriptionQuery = `
		INSERT INTO user_subscriptions (user_id, channel_id)
		VALUES ($1, $2)
		RETURNING id, user_id, channel_id, subscribed_at
	`

	getSubscriptionByIDQuery = `
		SELECT id, user_id, channel_id, subscribed_at
		FROM user_subscriptions
		WHERE id = $1
	`

	getSubscriptionByUserAndChannelQuary = `
		SELECT id, user_id, channel_id, subscribed_at
			FROM user_subscriptions
		WHERE user_id = $1 AND channel_id = $2
	`

	countAllUserSubscriptionsDQuery = `
		SELECT COUNT(id) AS total_count
			FROM user_subscriptions
		WHERE user_id = $1
	`

	getSubscriptionsByUserIDQuery = `
		SELECT us.id, us.user_id, us.channel_id, us.subscribed_at, c.id, c.username, c.title, c.description, c.parse_interval, c.last_parsed_at, c.created_at, c.updated_at
			FROM user_subscriptions us
		LEFT JOIN channels c ON us.channel_id = c.id
			WHERE user_id = $1
				ORDER BY subscribed_at DESC
			LIMIT $2 OFFSET $3
	`

	deleteSubscriptionQuery = `
		DELETE FROM user_subscriptions
		WHERE id = $1
		RETURNING id
	`
)
