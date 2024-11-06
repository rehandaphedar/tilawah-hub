package db

import (
	"database/sql"
	"log"

	"git.sr.ht/~rehandaphedar/tilawah-hub/internal/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

var Queries *sqlc.Queries

func Connect() {
	db, err := sql.Open("sqlite3", "file:data/db.sqlite?_fk=true")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	Queries = sqlc.New(db)
}
