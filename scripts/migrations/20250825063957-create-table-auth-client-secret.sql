CREATE TABLE IF NOT EXISTS auth_client_secret (
		client_id VARCHAR(32) NOT NULL,
		secret VARCHAR(256) NOT NULL,
		PRIMARY KEY(client_id, secret)
)