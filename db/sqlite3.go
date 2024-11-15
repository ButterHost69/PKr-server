package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
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

	workspaceTableQuery := `CREATE TABLE IF NOT EXISTS workspaces (
		username TEXT,
		workspace_name TEXT,

		PRIMARY KEY(username, workspace_name)
	);`

	currentUserIPTableQuery := `CREATE TABLE IF NOT EXISTS currentuserip (
		username TEXT PRIMARY KEY,
		ip_addr TEXT,
		port TEXT
	);`

	_, err = tx.Exec(usersTableQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(workspaceTableQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(currentUserIPTableQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		if rollback_err := tx.Rollback(); rollback_err != nil {
			return fmt.Errorf("could Not RollBack transaction during a commit error.\nError: %v", rollback_err)
		}
		return fmt.Errorf("could not Commit transaction.\nError: %v", err)
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
		return fmt.Errorf("error: Could Not Create Tables.\nError: %v", err)
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

// Returns Bool, if bool=false and err=nil, username or password incorrect
func RegisterNewWorkspace(username, password, workspace_name string) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	ifAuth, err := authUser(tx, username, password)
	if err != nil {
		return false, fmt.Errorf("error Could not Auth User.\nError: %v", err)
	}

	if !ifAuth {
		return false, nil
	}

	query := "INSERT INTO workspaces (username, workspace_name) VALUES (?,?)"
	if _, err = tx.Exec(query, username, workspace_name); err != nil {
		tx.Rollback()
		return false, fmt.Errorf("error Could not Execute Insert Statement for Register Workspace.\nError: %v", err)
	}

	if err = tx.Commit(); err != nil {
		if rollback_err := tx.Rollback(); rollback_err != nil {
			return false, fmt.Errorf("could Not RollBack transaction during a commit error.\nError: %v", rollback_err)
		}
		return false, fmt.Errorf("could not Commit transaction.\nError: %v", err)
	}

	return true, nil
}

func authUser(tx *sql.Tx, username, password string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE username=? AND password=?"
	rows, err := tx.Query(query, username, password)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	defer rows.Close()
	
	// Check if any rows retrieved
	if !rows.Next() {
		return false, nil
	}

	return true, nil
}


func UpdateUserIP(username, password, ip_addr, port string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	ifAuth, err := authUser(tx, username, password)
	if err != nil {
		return fmt.Errorf("error Could not Auth User.\nError: %v", err)
	}

	if !ifAuth {
		return fmt.Errorf("error Incorrect user credentials.\nError: %v", err)
	}

	// query := `UPDATE TABLE currentuserip 
	// SET ip_addr=?, port=?
	// WHERE username=?`

	query := `INSERT OR REPLACE INTO currentuserip (username, ip_addr, port) 
	VALUES (?,?,?);`

	_, err = tx.Exec(query, username, ip_addr, port)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error Could not Update Users IP.\nError: %v", err)
	}

	if err = tx.Commit(); err != nil {
		if rollback_err := tx.Rollback(); rollback_err != nil {
			return fmt.Errorf("could Not RollBack transaction during a commit error.\nError: %v", rollback_err)
		}
		return fmt.Errorf("could not Commit transaction.\nError: %v", err)
	}
	
	return nil
}

func CloseSQLiteDatabase() {
	db.Close()
}
