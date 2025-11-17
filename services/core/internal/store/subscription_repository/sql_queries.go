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

	getAllSubscriptionsQuery = `
		SELECT id, user_id, channel_id, subscribed_at
		FROM user_subscriptions
		ORDER BY subscribed_at DESC
		LIMIT $1 OFFSET $2
	`

	getSubscriptionsByUserIDQuery = `
		SELECT id, user_id, channel_id, subscribed_at
		FROM user_subscriptions
		WHERE user_id = $1
		ORDER BY subscribed_at DESC
		LIMIT $2 OFFSET $3
	`

	updateSubscriptionQuery = `
		UPDATE user_subscriptions
		SET user_id = $2, channel_id = $3
		WHERE id = $1
		RETURNING id, user_id, channel_id, subscribed_at
	`

	deleteSubscriptionQuery = `
		DELETE FROM user_subscriptions
		WHERE id = $1
		RETURNING id
	`
)
