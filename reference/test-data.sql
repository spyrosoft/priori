\i database.sql

INSERT INTO users (email_hash, password_hash, password_salt) VALUES ('\x6a67524f5d3117d9481dd39fdcffcde682b262e8ebbf64a81a26207413b178d6', '\x228640fa2cdf5cbf6a3e7964ac3035e59a62ddd113a4729b0d2057f2f79e703e', '\x305f94a6b2fb3de2897d');

--INSERT INTO user_tasks (user_id, task, short_term, long_term, urgency, difficulty) VALUES (1, 'My Projects', 10, 10, 0, 4);


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