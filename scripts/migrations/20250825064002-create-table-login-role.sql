CREATE TABLE IF NOT EXISTS login_role (
		login_id UUID NOT NULL,
		role_id VARCHAR(12) NOT NULL,
		PRIMARY KEY(login_id, role_id)
)