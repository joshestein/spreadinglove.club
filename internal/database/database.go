package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"spreadlove/db"
	sqlFiles "spreadlove/sql"
)

// Opens a connection to the SQLite database and initializes the schema.
// It returns both the raw *sql.DB and the sqlc-generated *db.Queries for type-safe operations.
func Setup(dbPath string) (*sql.DB, *db.Queries, error) {
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, nil, err
	}

	queries := db.New(database)

	// Setup DB from schema.sql
	_, err = database.Exec(sqlFiles.SchemaSQL)
	if err != nil {
		database.Close()
		return nil, nil, err
	}

	return database, queries, nil
}
