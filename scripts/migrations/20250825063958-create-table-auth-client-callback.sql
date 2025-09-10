CREATE TABLE IF NOT EXISTS auth_client_callback (
		client_id VARCHAR(32) NOT NULL,
		callback_url VARCHAR(256) NOT NULL,
		PRIMARY KEY(client_id, callback_url)
)