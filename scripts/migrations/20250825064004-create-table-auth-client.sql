CREATE TABLE IF NOT EXISTS auth_client (
		id VARCHAR(32) NOT NULL,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(1000),
		homepage_url VARCHAR(500) NOT NULL,
		callback_url VARCHAR(500) NOT NULL,
		PRIMARY KEY(id)
)