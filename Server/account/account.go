package account

import (
	"botify/main/db"
	"botify/main/logger"
	"botify/main/script"
	"os"
	"strconv"
)

type Account struct {
	id      int
	user    User
	scripts []script.Script
}

// ------ Account functions ------

func (a *Account) GetId() int {
	return a.id
}

func (a *Account) GetUser() *User {
	return &a.user
}

func (a *Account) GetScripts() []script.Script {
	return a.scripts
}

func (a *Account) GetScriptNames() []string {
	var names []string
	for _, s := range a.scripts {
		names = append(names, s.GetName())
	}
	return names
}

func CreateAccount(user User, db db.Database) Account {
	var a Account
	a.user = user
	a.id = user.GetId()

	ids := db.GetScriptIdsFromUid(a.id)

	for _, id := range ids {
		id, uid, name, description, path := db.GetScriptData(id)
		s := script.CreateScript(uid, id, name, description, path)
		a.scripts = append(a.scripts, s)
	}

	return a
}

// ------ Script functions ------

func (a *Account) AddScript(name string, description string, path string, db db.Database) {
	id := db.AddScript(a.id, name, description, path)
	s := script.CreateScript(a.id, id, name, description, path)
	a.scripts = append(a.scripts, s)
}

func (a *Account) DeleteScript(id int, db db.Database) {
	name := db.GetFieldById(id, "scripts", "name")
	strId := strconv.FormatInt(int64(id), 10)

	logger.Info("deleting script: " + name + " with id: " + strId + " from user: " + a.user.GetUsername())
	var path string = a.user.GetUsername() + "/scripts/" + name + ".lua"
	os.Remove("userdata/" + path)

	db.DeleteScript(id)

	var scripts []script.Script
	for _, s := range a.scripts {
		if s.GetId() != id {
			scripts = append(scripts, s)
		}
	}
	a.scripts = scripts
}
