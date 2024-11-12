package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db 	*sql.DB
)

func createAllTables() error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error: Could Not Initiate the Transaction.\nError: %v", err)
	}

	// TODO: [ ] Send Hashed Passwords from the client side
	// TODO : [ ] Ideally server should recv encrypted passwords (IDK How ??)
	usersTableQuery := `CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT
	);`

	_, err = tx.Exec(usersTableQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err!=nil {
		return err
	}

	return nil

}

func InitSQLiteDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./server_database.db")
	if err != nil {
		return fmt.Errorf("error: Could Not Start The Database.\nError: %v", err)
	}

	err = createAllTables()
	if err != nil {
		return fmt.Errorf("error: Could Create Tables.\nError: %v", err)
	}
	return nil
}

func CreateNewUser(username, password string) error {
	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err := db.Exec(query, username, password)
	if err != nil {
		return fmt.Errorf("error: Could Create New User %s .\nError: %v", username, err)
	}
	
	return nil
}

func CloseSQLiteDatabase(){
	db.Close()
}