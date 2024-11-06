CREATE TABLE sessions (
    session_token TEXT PRIMARY KEY NOT NULL,
    csrf_token TEXT NOT NULL,
	username VARCHAR(64) NOT NULL REFERENCES users (username) ON DELETE CASCADE 
);
