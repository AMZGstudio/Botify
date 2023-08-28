package script

import (
	"botify/main/db"
	"os"
)

// so the script is just a lua file which is kept in the userdata/username/scripts folder
// the script is named by the id of the script in the database
// the script is also kept in the database an a table called scripts

type Script struct {
	id          int
	uid         int
	name        string
	description string
	path        string
}

func (s *Script) GetId() int {
	return s.id
}

func (s *Script) GetUid() int {
	return s.uid
}

func (s *Script) GetName() string {
	return s.name
}

func (s *Script) GetDescription() string {
	return s.description
}

func (s *Script) GetPath() string {
	return s.path
}

func CreateScript(uid int, id int, name string, description string, path string) Script {
	var s Script
	s.uid = uid
	s.id = id
	s.name = name
	s.description = description
	s.path = path
	return s
}

func (s *Script) DeleteScript(id int, db db.Database) {
	db.DeleteScript(id)
	var username string = db.GetFieldById(s.uid, "user", "name")
	os.Remove("userdata/" + username + "/scripts/" + s.name + ".lua")
}
