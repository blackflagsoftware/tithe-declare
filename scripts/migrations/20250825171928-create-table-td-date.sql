CREATE TABLE IF NOT EXISTS td_date (
	id INT PRIMARY KEY,
	date_value DATE NOT NULL UNIQUE,
	hold DATE,
	confirm DATE,
	name VARCHAR(255),
	phone VARCHAR(255),
	email VARCHAR(255)
)