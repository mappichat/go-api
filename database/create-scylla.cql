-- Run this file on your database to setup the schemas
CREATE KEYSPACE IF NOT EXISTS user_data WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS user_data.posts (
	id uuid PRIMARY KEY,
	title text,
	body text,
	user_handle text,
	timestamp timestamp,
	latitude float,
	longitude float,
	level int,
	reply_count int,
	up_votes int,
	down_votes int,
);

CREATE TABLE IF NOT EXISTS user_data.replies (
	post_id uuid,
	id uuid,
	body text,
	user_handle text,
	timestamp timestamp,
	PRIMARY KEY (post_id, id),
);

CREATE TABLE IF NOT EXISTS user_data.votes (
	post_id uuid,
	user_handle text,
	up boolean,
	timestamp timestamp,
	PRIMARY KEY (post_id, user_handle),
);
