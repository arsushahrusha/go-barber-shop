package db

const (
	registerUserQuery = `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING id, login, password, created_at
	`

	getUserByLoginQuery = `
		SELECT id, login, password, created_at
		FROM users
		WHERE login = $1
	`

	createSessionQuery = `
		INSERT INTO sessions (user_id, expires_at)
		VALUES ($1, NOW() + INTERVAL '24 hours')
		RETURNING session_id, user_id, created_at, expires_at
	`
)