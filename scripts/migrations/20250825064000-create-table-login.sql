CREATE TABLE IF NOT EXISTS login (
		id UUID NOT NULL,
		email_addr VARCHAR(100) NOT NULL,
		first_name VARCHAR(50),
		last_name VARCHAR(100),
		pwd VARCHAR(250) NOT NULL,
		active BOOL DEFAULT true NOT NULL,
		set_pwd BOOL DEFAULT false NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NULL,
		PRIMARY KEY(id)
)