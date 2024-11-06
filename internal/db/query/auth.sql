-- name: AuthInsertUser :one
INSERT INTO users (username, password, displayname)
    VALUES (?1, ?2, ?3)
RETURNING
    username, displayname;

-- name: AuthSelectUser :one
SELECT
    username, password
FROM
    users
WHERE
    username = ?1;

-- name: AuthInsertSession :one
INSERT INTO sessions (session_token, csrf_token, username)
	VALUES (?1, ?2, ?3)
RETURNING
	*;

-- name: AuthSelectSession :one
SELECT
	*
FROM
	sessions
WHERE
	session_token = ?1;

-- name: AuthDeleteSession :one
DELETE
FROM
	sessions
WHERE
	session_token = ?1
RETURNING
	*;
