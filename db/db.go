package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase(databaseName string) error {
	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		return err
	}
	DB = db
	_, errCities := DB.Exec(
		`CREATE TABLE IF NOT EXISTS cities (
			id INTEGER PRIMARY KEY,
			name TEXT
		)`,
	)
	if errCities != nil {
		return errCities
	}
	_, errBorders := DB.Exec(
		"CREATE TABLE IF NOT EXISTS borders (`from` INTEGER, `to` INTEGER, UNIQUE(`from`, `to`))",
	)
	return errBorders
}
