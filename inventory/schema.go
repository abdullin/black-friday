package inventory

import "database/sql"

import _ "embed"

//go:embed schema.sql
var schema string

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
