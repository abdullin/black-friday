package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"text/template"
)

import _ "embed"

//go:embed schema.sql.tmpl
var schema string

type Config struct {
	Strict bool
}

func CreateSchema(db *sql.DB, strict bool) error {

	tpl, err := template.New("schema").Parse(schema)
	if err != nil {
		log.Fatalln(err)
	}
	buf := new(bytes.Buffer)

	err = tpl.Execute(buf, Config{Strict: strict})
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(buf.String())
	if err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}
