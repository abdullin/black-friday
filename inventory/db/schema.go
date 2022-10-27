package db

import (
	"database/sql"
	"fmt"
)

import _ "embed"

//go:embed schema.sql
var schema string

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}
