package db

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewSQLite(path string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	return db, nil
}
