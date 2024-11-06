CREATE TABLE recitations(
	 slug VARCHAR(64) NOT NULL,
	 name VARCHAR(64) NOT NULL,
	 reciter VARCHAR(64) NOT NULL REFERENCES users (username) ON DELETE CASCADE,
	 PRIMARY KEY(reciter, slug)
);
