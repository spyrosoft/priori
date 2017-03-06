DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS user_tasks;
DROP TABLE IF EXISTS users;


CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	email_hash BYTEA UNIQUE NOT NULL,
	password_hash BYTEA NOT NULL,
	password_salt BYTEA NOT NULL
);


CREATE TABLE user_sessions (
	user_id INT REFERENCES users(id),
	token TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT 'now'
);


CREATE TABLE user_tasks (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	parent_id INT,
	task TEXT NOT NULL,
	short_term INT NOT NULL,
	long_term INT NOT NULL,
	urgency INT NOT NULL,
	difficulty INT NOT NULL
);
