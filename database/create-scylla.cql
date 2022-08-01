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
)