package db

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"time"
)

type PGConfig struct {
	Host         string
	Port         string
	User         string
	Pass         string
	DB           string
	SslSkip      bool
	Timeout      int
	DialTimeout  int
	ReadTimeout  int
	WriteTimeout int
	QueryDebug   bool
}

func MakePGConnect(cfg PGConfig) (*bun.DB, error) {
	opts := []pgdriver.Option{
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		pgdriver.WithUser(cfg.User),
		pgdriver.WithPassword(cfg.Pass),
		pgdriver.WithDatabase(cfg.DB),
	}
	if cfg.SslSkip {
		opts = append(opts, pgdriver.WithInsecure(true))
	}
	if cfg.Timeout != 0 {
		opts = append(opts, pgdriver.WithTimeout(time.Duration(cfg.Timeout)*time.Second))
	}
	
	if cfg.DialTimeout != 0 {
		opts = append(opts, pgdriver.WithDialTimeout(time.Duration(cfg.DialTimeout)*time.Second))
	}
	
	if cfg.ReadTimeout != 0 {
		opts = append(opts, pgdriver.WithReadTimeout(time.Duration(cfg.ReadTimeout)*time.Second))
	}
	
	if cfg.WriteTimeout != 0 {
		opts = append(opts, pgdriver.WithWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Second))
	}
	pgConn := pgdriver.NewConnector(opts...)
	sqlPG := sql.OpenDB(pgConn)
	database := bun.NewDB(sqlPG, pgdialect.New())
	if cfg.QueryDebug {
		database.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}
	return database, nil
}
