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

	createOrderQuery = `
		INSERT INTO orders (user_id, status)
		VALUES ($1, 'CREATED')
		RETURNING id, user_id, status, created_at, updated_at
	`

	getOrdersByUserIDQuery = `
		SELECT id, user_id, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	getActiveOrdersByUserIDQuery = `
		SELECT id, user_id, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1 AND status NOT IN ('DONE', 'CANCELLED')
		ORDER BY created_at DESC
	`

	getSessionByUserIDQuery = `
		SELECT session_id, user_id, created_at, expires_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY created_At DESC
		LiMIT 1
	`

	updateSessionExpiryQuery = `
		UPDATE sessions
		SET expires_at = NOW() + INTERVAL '24 hours'
		WHERE session_id = $1
	`

	changeOrderStatusQuery = `
		UPDATE orders
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

)