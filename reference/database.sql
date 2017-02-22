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



INSERT INTO users (email_hash, password_hash, password_salt) VALUES ('\x6a67524f5d3117d9481dd39fdcffcde682b262e8ebbf64a81a26207413b178d6', '\x228640fa2cdf5cbf6a3e7964ac3035e59a62ddd113a4729b0d2057f2f79e703e', '\x305f94a6b2fb3de2897d');
INSERT INTO user_tasks (user_id, task, short_term, long_term, urgency, difficulty) VALUES (1, 'My Projects', 10, 10, 0, 4);
INSERT INTO user_tasks (user_id, parent_id, task, short_term, long_term, urgency, difficulty) VALUES (1, 1, 'My Projects Child', 4, 4, 0, 4);


-- with recursive all_posts (id, parentid, root_id) as (
--     select t1.id,
--     t1.parent_forum_post_id as parentid,
--     t1.id as root_id
--     from forumposts t1

--     union all

--     select c1.id,
--     c1.parent_forum_post_id as parentid,
--     p.root_id
--     from forumposts c1
--     join all_posts p on p.id = c1.parent_forum_post_id
-- )
-- DELETE FROM forumposts
--  WHERE id IN (SELECT id FROM all_posts WHERE root_id=1349);
-- Source: http://stackoverflow.com/questions/10381243/delete-recursive-children#10381384