CREATE TABLE IF NOT EXISTS accounts (
  id UUID NOT NULL,
	email VARCHAR NOT NULL UNIQUE,
  user_handle VARCHAR NOT NULL UNIQUE,
	constraint check_lowercase_user_handle check (lower(user_handle) = user_handle),
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS posts (
	id UUID NOT NULL,
  account_id UUID NOT NULL REFERENCES accounts(id),
  title VARCHAR NOT NULL,
	body VARCHAR NOT NULL,
  post_level SMALLINT NOT NULL,
	latitude REAL NOT NULL,
	longitude REAL NOT NULL,
  time_stamp TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS replies (
	id UUID NOT NULL,
  post_id UUID NOT NULL REFERENCES posts(id),
  account_id UUID NOT NULL REFERENCES accounts(id),
	body VARCHAR NOT NULL,
	latitude REAL NOT NULL,
	longitude REAL NOT NULL,
	time_stamp TIMESTAMP WITH TIME ZONE NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS votes (
	post_id UUID NOT NULL REFERENCES posts(id),
	account_id UUID NOT NULL REFERENCES accounts(id),
	vote_weight REAL NOT NULL,
	vote_level int,
	latitude float,
	longitude float,
	time_stamp timestamp,
  PRIMARY KEY (post_id, account_id, vote_level)
);
