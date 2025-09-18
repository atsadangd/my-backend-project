package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Init opens (and creates if needed) the sqlite database file and ensures the users table exists.
func Init(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	// simple ping to validate
	if err := DB.Ping(); err != nil {
		return err
	}

	create := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		first_name TEXT,
		last_name TEXT,
		phone TEXT,
		avatar TEXT
	);`
	if _, err := DB.Exec(create); err != nil {
		return err
	}
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
