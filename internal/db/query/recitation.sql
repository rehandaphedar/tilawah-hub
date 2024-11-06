-- name: RecitationCreateRecitation :one
INSERT INTO recitations(reciter, slug, name)
	VALUES (?1, ?2, ?3)
RETURNING *;

-- name: RecitationSelectRecitations :many
SELECT
	*
FROM
    recitations;

-- name: RecitationSelectRecitation :one
SELECT
	*
FROM
    recitations
WHERE
	reciter = ?1 AND slug = ?2;

-- name: RecitationUpdateRecitation :one
UPDATE recitations
SET
	name = ?3
WHERE
	reciter = ?1 AND slug = ?2
RETURNING *;

-- name: RecitationDeleteRecitation :one
DELETE FROM recitations
WHERE
	reciter = ?1 AND slug = ?2
RETURNING *;
