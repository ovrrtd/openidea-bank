package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	DB_HOST     = os.Getenv("DB_HOST")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_PARAMS   = os.Getenv("DB_PARAMS")
)

func NewDBDefaultSql() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s",
		DB_USERNAME,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
		DB_PARAMS,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = prometheus.Register(collectors.NewDBStatsCollector(db, DB_NAME))
	if err != nil {
		return nil, err
	}

	return db, nil
}
