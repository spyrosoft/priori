DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS list_tasks;

DROP TABLE IF EXISTS user_lists;

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


CREATE TABLE user_lists (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	name TEXT NOT NULL
);

CREATE TABLE list_tasks (
	id SERIAL PRIMARY KEY,
	list_id INT REFERENCES user_lists(id),
	title TEXT NOT NULL,
	difficulty INT NOT NULL,
	short_term INT NOT NULL,
	long_term INT NOT NULL
);

INSERT INTO users (email_hash, password_hash, password_salt) VALUES ('\x6a67524f5d3117d9481dd39fdcffcde682b262e8ebbf64a81a26207413b178d6', '\x228640fa2cdf5cbf6a3e7964ac3035e59a62ddd113a4729b0d2057f2f79e703e', '\x305f94a6b2fb3de2897d');
