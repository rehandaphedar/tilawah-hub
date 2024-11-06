// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: auth.sql

package sqlc

import (
	"context"
)

const authDeleteSession = `-- name: AuthDeleteSession :one
DELETE
FROM
	sessions
WHERE
	session_token = ?1
RETURNING
	session_token, csrf_token, username
`

func (q *Queries) AuthDeleteSession(ctx context.Context, sessionToken string) (Session, error) {
	row := q.db.QueryRowContext(ctx, authDeleteSession, sessionToken)
	var i Session
	err := row.Scan(&i.SessionToken, &i.CsrfToken, &i.Username)
	return i, err
}

const authInsertSession = `-- name: AuthInsertSession :one
INSERT INTO sessions (session_token, csrf_token, username)
	VALUES (?1, ?2, ?3)
RETURNING
	session_token, csrf_token, username
`

type AuthInsertSessionParams struct {
	SessionToken string `json:"session_token"`
	CsrfToken    string `json:"csrf_token"`
	Username     string `json:"username"`
}

func (q *Queries) AuthInsertSession(ctx context.Context, arg AuthInsertSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, authInsertSession, arg.SessionToken, arg.CsrfToken, arg.Username)
	var i Session
	err := row.Scan(&i.SessionToken, &i.CsrfToken, &i.Username)
	return i, err
}

const authInsertUser = `-- name: AuthInsertUser :one
INSERT INTO users (username, password, displayname)
    VALUES (?1, ?2, ?3)
RETURNING
    username, displayname
`

type AuthInsertUserParams struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Displayname string `json:"displayname"`
}

type AuthInsertUserRow struct {
	Username    string `json:"username"`
	Displayname string `json:"displayname"`
}

func (q *Queries) AuthInsertUser(ctx context.Context, arg AuthInsertUserParams) (AuthInsertUserRow, error) {
	row := q.db.QueryRowContext(ctx, authInsertUser, arg.Username, arg.Password, arg.Displayname)
	var i AuthInsertUserRow
	err := row.Scan(&i.Username, &i.Displayname)
	return i, err
}

const authSelectSession = `-- name: AuthSelectSession :one
SELECT
	session_token, csrf_token, username
FROM
	sessions
WHERE
	session_token = ?1
`

func (q *Queries) AuthSelectSession(ctx context.Context, sessionToken string) (Session, error) {
	row := q.db.QueryRowContext(ctx, authSelectSession, sessionToken)
	var i Session
	err := row.Scan(&i.SessionToken, &i.CsrfToken, &i.Username)
	return i, err
}

const authSelectUser = `-- name: AuthSelectUser :one
SELECT
    username, password
FROM
    users
WHERE
    username = ?1
`

type AuthSelectUserRow struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) AuthSelectUser(ctx context.Context, username string) (AuthSelectUserRow, error) {
	row := q.db.QueryRowContext(ctx, authSelectUser, username)
	var i AuthSelectUserRow
	err := row.Scan(&i.Username, &i.Password)
	return i, err
}