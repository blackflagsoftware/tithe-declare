CREATE TABLE IF NOT EXISTS auth_authorize (
		id VARCHAR(32) NOT NULL,
		client_id VARCHAR(32) NOT NULL,
		verifier TEXT,
		verifier_encode_method VARCHAR(10),
		state VARCHAR(100),
		scope VARCHAR(256),
		authorized_at TIMESTAMP NOT NULL,
		auth_code_at TIMESTAMP,
		auth_code VARCHAR(256),
		PRIMARY KEY(id)
)