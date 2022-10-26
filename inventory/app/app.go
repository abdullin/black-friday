package app

import (
	"database/sql"
)

// App is used by features to implement their functionality
type App struct {
	db *sql.DB
}

func New(db *sql.DB) *App {
	if db == nil {
		panic("db is nil")
	}
	return &App{
		db: db,
	}
}
