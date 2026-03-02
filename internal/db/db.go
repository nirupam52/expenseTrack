package db

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

func OpenDB(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath);
}
