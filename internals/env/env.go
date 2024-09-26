package env

import (
	"log/slog"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/robrohan/laffaire/internals/models"
	"github.com/robrohan/laffaire/internals/repository"
)

// Env context for db, logger, etc. This is passed within a request
type Env struct {
	Db        *sqlx.DB
	Log       *slog.Logger
	Cfg       *models.Config
	Router    *mux.Router
	User      *models.User
	RandState string
	Repo      *repository.DataRepository
}
