package userman

import (
	"botify/main/db"
	"crypto/sha256"
	"encoding/hex"
)

type User struct {
	id        int
	username  string // unique
	password  string // hash
	email     string // if null, no email
	phone     string // if null, no phone
	confirmed bool
}

func (u *User) Login(username string, password string, db db.Database) {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	u.username = username
	u.password = hashedPassword
	u.confirmed = db.IsValidUser(u.username, u.password)
	u.id = db.GetIdFromUsername(u.username)

	// if id is -1, the user was not found
	if u.id == -1 {
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
func (u *User) Signup(username string, password string, email string, phone string, db db.Database) {
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
	u.id = db.AddUser(u.username, u.password, u.email, u.phone)

	// if id is -1, the user was not added
	if u.id == -1 {
		u.confirmed = false
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
