package repo

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
	"server/config"
)

func NewDBConnect(c *config.Credentials) (*sqlx.DB, error) {
	dbConf, err := pgx.ParseConfig(c.GetPostgreDBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse config")
	}

	dbConf.Host = c.PostgreSQL.Host
	logger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.JSONFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.ErrorLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	dbConf.Logger = logrusadapter.NewLogger(logger)
	dsn := stdlib.RegisterConnConfig(dbConf)

	sql.Register("wrapper", stdlib.GetDefaultDriver())
	wdb, err := sql.Open("wrapper", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database")
	}

	const (
		maxOpenConns    = 50
		maxIdleConns    = 50
		connMaxLifetime = 3
		connMaxIdleTime = 5
	)
	db := sqlx.NewDb(wdb, "pgx")
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database, unavailable resourse")
	}

	return db, nil
}
