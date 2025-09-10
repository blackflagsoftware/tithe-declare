CREATE TABLE IF NOT EXISTS auth_refresh (
		client_id VARCHAR(32) NOT NULL,
		token VARCHAR(256) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		PRIMARY KEY(client_id, token)
)