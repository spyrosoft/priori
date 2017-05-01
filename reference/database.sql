DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;


CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	email_hash BYTEA UNIQUE NOT NULL,
	password_hash BYTEA NOT NULL,
	password_salt BYTEA NOT NULL,
	-- Storing tasks as TEXT makes sense
	-- We won't be performing database actions on the JSON
	tasks TEXT
);


CREATE TABLE user_sessions (
	user_id INT REFERENCES users(id),
	token TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT 'now'
);
