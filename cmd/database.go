package main

import (
	"database/sql"
	"fmt"
)

func newDBDefaultSql() (*sql.DB, error) {
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

	return db, nil
}
