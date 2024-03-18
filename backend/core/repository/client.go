package repository

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
)

type Config struct {
	URL       string `env:"DATABASE_URL"`
	DebugMode string `env:"DB_DEBUG"`
}

func NewDatabase(cfg *Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	if cfg.DebugMode != "" {
		db.AddQueryHook(dbLogger{})
	}
	return db, nil
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	if query, err := q.FormattedQuery(); err != nil {
		utils.Logger.Debugf("[SQL]: %s", err.Error())
	} else {
		utils.Logger.Debugf("[SQL]: ", string(query))
	}

	return nil
}
