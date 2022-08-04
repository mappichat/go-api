-- Run this file on your database to setup the schemas
CREATE KEYSPACE IF NOT EXISTS user_data WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS user_data.posts (
	id uuid,
	part_key text,
	tile text,
	title text,
	body text,
	account_id text,
	user_handle text,
	time_stamp timestamp,
	latitude float,
	longitude float,
	level int,
	reply_count int,
	up_votes int,
	down_votes int,
	PRIMARY KEY ((part_key), level, tile, id, account_id),
);

CREATE TABLE IF NOT EXISTS user_data.replies (
	post_id uuid,
	id uuid,
	body text,
	account_id text,
	user_handle text,
	tile text,
	latitude float,
	longitude float,
	time_stamp timestamp,
	PRIMARY KEY ((post_id), id, account_id),
);

CREATE TABLE IF NOT EXISTS user_data.votes (
	post_id uuid,
	account_id text,
	up boolean,
	level int,
	tile text,
	latitude float,
	longitude float,
	time_stamp timestamp,
	PRIMARY KEY ((post_id), level, account_id),
);
