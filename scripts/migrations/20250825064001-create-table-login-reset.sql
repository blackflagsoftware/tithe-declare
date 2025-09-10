CREATE TABLE IF NOT EXISTS login_reset (
		login_id UUID NOT NULL,
		reset_token UUID NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NULL,
		PRIMARY KEY(login_id, reset_token)
)