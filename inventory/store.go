package inventory

import "database/sql"

import _ "embed"

//go:embed schema.sql
var schema string

type Store struct {
	db *sql.DB
}

func CreateSchema(db *sql.DB) error {
	_, err := db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return err
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	} else {
		defer tx.Rollback()
	}

	_, err = tx.Exec(schema)
	if err != nil {
		return err
	}

	return tx.Commit()

}
