-- name: UserSelectUsers :many
SELECT
    username, displayname
FROM
    users;

-- name: UserSelectUser :one
SELECT
    username, displayname
FROM
    users
WHERE
    username = ?1;

-- name: UserUpdateUser :one
UPDATE users
SET
	displayname = ?2
WHERE
	username = ?1
RETURNING
	username,
	displayname;

-- name: UserDeleteUser :one
DELETE FROM
	users
WHERE
	username = ?1
RETURNING
	username,
	displayname;
