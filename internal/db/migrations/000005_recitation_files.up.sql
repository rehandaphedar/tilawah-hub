CREATE TABLE recitation_files(
	 reciter VARCHAR(64) NOT NULL,
	 slug VARCHAR(64) NOT NULL,
	 verse_key VARCHAR(6) NOT NULL,
	 PRIMARY KEY(reciter, slug, verse_key),
	 FOREIGN KEY (reciter, slug) REFERENCES recitations(reciter, slug) ON DELETE CASCADE
);
