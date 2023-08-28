package db

import (
	"botify/main/logger"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	// the database
	db   *sql.DB
	path string
	open bool
}

// ------ Database functions ------

func (d *Database) SetPath(path string) {
	d.path = path
}

func (d *Database) Open() {
	if _, err := os.Stat(d.path); os.IsNotExist(err) {
		d.Create(d.path)
	}

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
	var err error

	d.db, err = sql.Open("sqlite3", path)
	if err != nil {
		logger.Error(err)
	}
	_, err = d.db.Exec("CREATE TABLE user (id INTEGER PRIMARY KEY, name TEXT UNIQUE, password TEXT, email TEXT, phone TEXT)")
	d.open = true
	if err != nil {
		logger.Error(err)
	}
	_, err = d.db.Exec("CREATE TABLE scripts (id INTEGER PRIMARY KEY, uid INTEGER, name TEXT, description TEXT, path TEXT, FOREIGN KEY(uid) REFERENCES user(id))")
	if err != nil {
		logger.Error(err)
	}
	// close the database
	d.db.Close()
	d.open = false
}

func (d *Database) GetIdFromName(name string, table string) int {
	if !d.open {
		d.Open()
	}

	// get the id from the name
	rows, err := d.db.Query(fmt.Sprintf("SELECT id FROM %s WHERE name = ?", table), name)
	if err != nil {
		logger.Error(err)
		return -1
	}

	// check if the id exists
	var id int = -1
	exists := rows.Next()
	if exists {
		rows.Scan(&id)
	}
	rows.Close()
	d.Close()
	return id
}

// get the name from the id
func (d *Database) GetFieldById(id int, field string, table string) string {
	if !d.open {
		d.Open()
	}

	// get the name from the id
	rows, err := d.db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", field, table), id)
	if err != nil {
		logger.Error(err)
		return ""
	}

	// check if the id exists
	var name string = ""
	exists := rows.Next()
	if exists {
		rows.Scan(&name)
	}
	rows.Close()
	d.Close()
	return name
}

// get all the ids from a table
func (d *Database) GetAllIds(table string) []int {
	if !d.open {
		d.Open()
	}

	rows, err := d.db.Query(fmt.Sprintf("SELECT id FROM %s", table))
	if err != nil {
		logger.Error(err)
	}

	d.Close()
	var ids []int
	var id int
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}

	rows.Close()
	return ids
}

// ------ User functions ------

func (d *Database) AddUser(username string, password string, email string, phone string) int {
	// check if the database is open
	if !d.open {
		d.Open()
	}

	// add the user to the database
	_, err := d.db.Exec("INSERT INTO user (name, password, email, phone) VALUES (?, ?, ?, ?)", username, password, email, phone)
	if err != nil {
		logger.Error(err)
	}

	// close the database
	d.Close()
	return d.GetIdFromName(username, "user")
}

func (d *Database) IsValidUser(username string, password string) bool {
	// check if the database is open
	if !d.open {
		d.Open()
	}

	// check if the user is valid
	rows, err := d.db.Query("SELECT * FROM user WHERE name = ? AND password = ?", username, password)
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
	rows, err := d.db.Query("SELECT * FROM user WHERE name = ?", username)
	if err != nil {
		logger.Error(err)
	}

	d.Close()
	exists := rows.Next()
	rows.Close()
	return !exists
}

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

// ------ Script functions ------

func (d *Database) AddScript(uid int, name string, description string, path string) int {
	if !d.open {
		d.Open()
	}

	_, err := d.db.Exec("INSERT INTO scripts (uid, name, description, path) VALUES (?, ?, ?, ?)", uid, name, description, path)
	if err != nil {
		panic(err)
	}

	d.Close()
	return d.GetIdFromName(name, "scripts")
}

func (d *Database) GetScriptIdsFromUid(uid int) []int {
	if !d.open {
		d.Open()
	}

	rows, err := d.db.Query("SELECT id FROM scripts WHERE uid = ?", uid)
	if err != nil {
		logger.Error(err)
	}

	d.Close()
	var ids []int
	var id int
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}

	rows.Close()
	return ids
}

func (d *Database) GetScriptData(id int) (int, int, string, string, string) {
	if !d.open {
		d.Open()
	}

	rows, err := d.db.Query("SELECT * FROM scripts WHERE id = ?", id)
	if err != nil {
		logger.Error(err)
	}

	d.Close()
	var uid int
	var name string
	var description string
	var path string
	for rows.Next() {
		rows.Scan(&id, &uid, &name, &description, &path)
	}

	rows.Close()
	return id, uid, name, description, path
}

func (d *Database) DeleteScript(id int) {
	logger.Info("deleting script with id: " + fmt.Sprint(id))
	if !d.open {
		d.Open()
	}

	_, err := d.db.Exec("DELETE FROM scripts WHERE id = ?", id)
	if err != nil {
		panic(err)
	}

	d.Close()
}
