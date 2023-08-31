package account

import (
	"botify/main/db"
	"botify/main/logger"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

type User struct {
	id             int
	username       string // unique
	password       string // hash
	email          string // if null, no email
	phone          string // if null, no phone
	confirmed      bool
	connectionType string // control or a service
}

func (u *User) Login(username string, password string, connectionType string, db db.Database) {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	u.username = username
	u.password = hashedPassword
	u.connectionType = connectionType
	u.confirmed = db.IsValidUser(u.username, u.password)
	// log all the data of the user
	u.id = db.GetIdFromName(u.username, "user")

	// if id is -1, the user was not found
	if u.id == -1 {
		logger.Warn("user: " + u.username + " with password: " + u.password + " was not found")
		u.confirmed = false
	}
}

func (u *User) Logout() {
	// reset everything
	u.username = ""
	u.password = ""
	u.confirmed = false
	u.id = -1
}

// create a user
func (u *User) Signup(username string, password string, email string, phone string, connectionType string, db db.Database) {
	// check if the username is available
	if !db.IsUsernameAvailable(username) {
		u.confirmed = false
		return
	}
	// hash the password
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	u.username = username
	u.password = hashedPassword
	u.email = email
	u.phone = phone
	u.confirmed = true
	u.connectionType = connectionType
	u.id = db.AddUser(u.username, u.password, u.email, u.phone)

	// if id is -1, the user was not added
	if u.id == -1 {
		u.confirmed = false
	}
}

// create folders for the user
func (u *User) CreateFolders() {
	// check if userdata exists
	if _, err := os.Stat("userdata"); os.IsNotExist(err) {
		os.Mkdir("userdata", 0777)
	}

	// check if userdata/<username> exists
	if _, err := os.Stat("userdata/" + u.username); os.IsNotExist(err) {
		os.Mkdir("userdata/"+u.username, 0777)
	}

	// check if userdata/<username>/scripts exists
	if _, err := os.Stat("userdata/" + u.username + "/scripts"); os.IsNotExist(err) {
		os.Mkdir("userdata/"+u.username+"/scripts", 0777)
	}
}

// delete a user
func (u *User) Delete(db db.Database) {
	// use the database
	db.DeleteUser(u.id)

	u.username = ""
	u.password = ""
	u.confirmed = false
	u.id = -1
}

// create getters for the user
func (u *User) GetId() int {
	return u.id
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPhone() string {
	return u.phone
}

func (u *User) IsConfirmed() bool {
	return u.confirmed
}

func (u *User) GetConnectionType() string {
	return u.connectionType
}

// get a user by id
func GetUserById(id int, db db.Database) User {
	var u User
	u.id = id
	u.username = db.GetFieldById(id, "name", "user")
	u.password = db.GetFieldById(id, "password", "user")
	u.email = db.GetFieldById(id, "email", "user")
	u.phone = db.GetFieldById(id, "phone", "user")
	// try to login the user
	u.Login(u.username, u.password, "", db)

	return u
}
