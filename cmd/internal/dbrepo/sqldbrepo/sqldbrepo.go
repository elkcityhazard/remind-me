package sqldbrepo

import "github.com/elkcityhazard/remind-me/cmd/internal/config"

type SQLDBRepo struct {
	Config *config.AppConfig
}

func NewSQLDBRepo(ac *config.AppConfig) *SQLDBRepo {
	return &SQLDBRepo{
		Config: ac,
	}
}
