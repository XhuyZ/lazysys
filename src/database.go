
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./lazysys.db")
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS services (
		name TEXT PRIMARY KEY,
		description TEXT
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getServiceDescription(db *sql.DB, serviceName string) (string, error) {
	var description string
	err := db.QueryRow("SELECT description FROM services WHERE name = ?", serviceName).Scan(&description)
	if err != nil {
		if err == sql.ErrNoRows {
			return "No description found for this service.", nil
		}
		return "", err
	}
	return description, nil
}

func updateServiceDescription(db *sql.DB, serviceName, description string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO services (name, description) VALUES (?, ?)", serviceName, description)
	return err
}
