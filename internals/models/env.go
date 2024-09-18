package models

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Env context for db, logger, etc. This is passed within a request
type Env struct {
	Db        *sqlx.DB
	Log       *slog.Logger
	Cfg       *Config
	Router    *mux.Router
	User      *User
	RandState string
}
