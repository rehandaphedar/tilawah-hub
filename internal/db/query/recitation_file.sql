-- name: RecitationFileCreateRecitationFile :one
INSERT INTO recitation_files(reciter, slug, verse_key)
	VALUES (?1, ?2, ?3)
RETURNING *;

-- name: RecitationFileSelectRecitationFiles :many
SELECT
	*
FROM
    recitation_files
WHERE
	reciter = ?1 AND slug = ?2;

-- name: RecitationFileSelectRecitationFile :one
SELECT
	*
FROM
    recitation_files
WHERE
	reciter = ?1 AND slug = ?2 AND verse_key = ?3;

-- name: RecitationFileUpdateRecitationFile :one
UPDATE recitation_files
SET
	has_timings = ?4,
	lafzize_processing = ?5
WHERE
	reciter = ?1 AND slug = ?2 AND verse_key = ?3
RETURNING *;

-- name: RecitationFileDeleteRecitationFile :one
DELETE FROM recitation_files
WHERE
	reciter = ?1 AND slug = ?2 AND verse_key = ?3
RETURNING *;
