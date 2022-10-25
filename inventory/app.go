package inventory

import (
	. "black-friday/api"
	"black-friday/fx"
	"context"
	"database/sql"
)

type App struct {
	db *sql.DB
	UnimplementedInventoryServiceServer
}

func (s *App) GetTx(ctx context.Context) *fx.Tx {
	return fx.Begin(ctx, s.db)
}

func NewApp(db *sql.DB) *App {
	if db == nil {
		panic("db is nil")
	}
	return &App{
		db: db,
	}
}
