package user_repository

const (
	createUserQuery = `
		INSERT INTO users (telegram_id, username, is_active, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, telegram_id, username, created_at, updated_at, is_active, role
	`

	getUserByIDQuery = `
		SELECT id, telegram_id, username, created_at, updated_at, is_active, role
		FROM users
		WHERE id = $1
	`

	getAllUsersQuery = `
		SELECT id, telegram_id, username, created_at, updated_at, is_active, role
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	updateUserQuery = `
		UPDATE users
		SET telegram_id = $2, username = $3, is_active = $4, role = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, telegram_id, username, created_at, updated_at, is_active, role
	`

	deleteUserQuery = `
		DELETE FROM users
		WHERE id = $1
		RETURNING id
	`
)
