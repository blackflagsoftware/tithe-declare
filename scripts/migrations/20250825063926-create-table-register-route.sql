CREATE TABLE IF NOT EXISTS register_route (
		raw_path VARCHAR(255) NOT NULL,
		transformed_path VARCHAR(255) NOT NULL,
		roles JSON,
		PRIMARY KEY(raw_path)
)