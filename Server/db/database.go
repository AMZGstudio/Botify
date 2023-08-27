package db

import (
	"botify/main/logger"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	// the database
	db   *sql.DB
	path string
	open bool
}

func (d *Database) SetPath(path string) {
	d.path = path
}

func (d *Database) Open() {
	var err error
	// open the database
	d.db, err = sql.Open("sqlite3", d.path)
	d.open = true
	if err != nil {
		logger.Error(err)
	}
}

func (d *Database) Close() {
	// close the database
	d.db.Close()
	d.open = false
}

func (d *Database) Create(path string) {
	d.path = path
	// if it doesn't, create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		d.db, err = sql.Open("sqlite3", path)
		if err != nil {
			logger.Error(err)
		}
		_, err = d.db.Exec("CREATE TABLE user (id INTEGER PRIMARY KEY, username TEXT UNIQUE, password TEXT, email TEXT, phone TEXT")
		d.open = true
		if err != nil {
			logger.Error(err)
		}
	} else {
		d.db, err = sql.Open("sqlite3", path)
		d.open = true
		if err != nil {
			logger.Error(err)
		}
	}
	// close the database
	d.db.Close()
	d.open = false
}

// a function that adds a user to the database
func (d *Database) AddUser(username string, password string, email string, phone string) int {
	// check if the database is open
	if !d.open {
		d.Open()
	}

	// add the user to the database
	_, err := d.db.Exec("INSERT INTO user (username, password, email, phone) VALUES (?, ?, ?, ?)", username, password, email, phone)
	if err != nil {
		logger.Error(err)
	}

	// close the database
	d.Close()
	return d.GetIdFromUsername(username)
}

// check if a user is valid
func (d *Database) IsValidUser(username string, password string) bool {
	// check if the database is open
	if !d.open {
		d.Open()
	}

	// check if the user is valid
	rows, err := d.db.Query("SELECT * FROM user WHERE username = ? AND password = ?", username, password)
	if err != nil {
		logger.Error(err)
	}

	// close the database
	d.Close()
	exists := rows.Next()
	rows.Close()
	return exists
}

func (d *Database) IsUsernameAvailable(username string) bool {
	if !d.open {
		d.Open()
	}

	// check if the user is valid
	rows, err := d.db.Query("SELECT * FROM user WHERE username = ?", username)
	if err != nil {
		logger.Error(err)
	}

	d.Close()
	exists := rows.Next()
	rows.Close()
	return !exists
}

// get id from username
func (d *Database) GetIdFromUsername(username string) int {
	// check if the database is open
	if !d.open {
		d.Open()
	}

	// get the id from the username
	rows, err := d.db.Query("SELECT id FROM user WHERE username = ?", username)
	if err != nil {
		logger.Error(err)
		return -1
	}

	var id int
	rows.Next()
	rows.Scan(&id)
	rows.Close()
	d.Close()
	return id
}

// delete a user
func (d *Database) DeleteUser(id int) {
	if !d.open {
		d.Open()
	}

	_, err := d.db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		panic(err)
	}

	d.Close()
}
