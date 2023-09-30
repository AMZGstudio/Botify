package main

import (
	"botify/main/account"
	"botify/main/db"
	"botify/main/logger"
	"botify/main/script"
	"io"
	"os"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// a function of a loop that runs every 5 seconds in a goroutine
func handleAccounts(database *db.Database) {
	var acc []account.Account
	var running []script.Script
	// every half second, check if any scripts need to be run
	for {
		updateAccounts(&acc, *database)
		for _, a := range acc {
			// print the username for debugging
			for _, s := range a.GetScripts() {
				var found bool
				for _, r := range running {
					if s.GetId() == r.GetId() {
						found = true
						break
					}
				}
				if !found {
					// run the script
					running = append(running, s)
					go runScript(s, a, &running, database)

				}
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// this function will add and remove accounts from a given list
func updateAccounts(accPtr *[]account.Account, db db.Database) {
	acc := *accPtr
	ids := db.GetAllIds("user")

	// check if any accounts need to be removed
	var reAdd []account.Account
	for _, a := range acc {
		for _, id := range ids {
			if a.GetId() == id {
				reAdd = append(reAdd, a)
				break
			}
		}
	}
	acc = reAdd
	// check if any accounts need to be added
	for _, id := range ids {
		var found bool
		for _, a := range acc {
			if a.GetId() == id {
				found = true
				break
			}
		}
		if !found {
			// add the account
			user := account.GetUserById(id, db)
			a := account.CreateAccount(user, db)
			acc = append(acc, a)
		}
	}

	// now remake all the accounts to make sure they are up to date
	for i, a := range acc {
		acc[i] = account.CreateAccount(*a.GetUser(), db)
	}

	*accPtr = acc
}

// this function will be embedded into the lua script and they will be able to create files with it
func luaFileCreate(L *lua.LState) int {
	// get the file name
	fileName := L.ToString(1)
	// get the file contents
	fileContents := L.ToString(2)

	// create the file in a poc folder (if it doesnt exist, create it)
	if _, err := os.Stat("poc"); os.IsNotExist(err) {
		os.Mkdir("poc", 0755)
	}
	file, err := os.Create("poc/" + fileName)
	if err != nil {
		logger.Error(err)
		return 0
	}
	defer file.Close()

	// write the contents to the file
	_, err = file.WriteString(fileContents)
	if err != nil {
		logger.Error(err)
		return 0
	}

	return 0
}

func runScript(sc script.Script, acc account.Account, running *[]script.Script, db *db.Database) {
	// print for debugging
	logger.Info("running script: " + sc.GetName() + " by " + acc.GetUser().GetUsername())
	defer removeScriptFiles(sc, db, acc)
	defer removeScript(sc, running)

	// open or create the log file
	// if this file does not exist, create it
	if _, err := os.Stat("userdata/" + acc.GetUser().GetUsername() + "/logs"); os.IsNotExist(err) {
		os.Mkdir("userdata/"+acc.GetUser().GetUsername()+"/logs", 0755)
	}
	logFile, err := os.OpenFile("userdata/"+acc.GetUser().GetUsername()+"/logs/"+sc.GetName()+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(err)
		return
	}
	// get the time now
	now := time.Now()
	for {
		username := acc.GetUser().GetUsername()

		logger.InfoFile("executing script: "+sc.GetName()+" by "+username, logFile)

		L := lua.NewState()
		defer L.Close()

		L.SetGlobal("username", lua.LString(username))
		L.SetGlobal("script_name", lua.LString(sc.GetName()))
		L.SetGlobal("print", L.NewFunction(func(L *lua.LState) int { io.WriteString(logFile, L.ToString(1)+"\n"); return 0 }))
		L.SetGlobal("file_create", L.NewFunction(luaFileCreate))

		if err := L.DoFile("userdata/" + username + "/scripts/" + sc.GetName() + ".lua"); err != nil {
			logger.Error(err)
			return
		}

		time.Sleep(500 * time.Millisecond)

		// check if it has been 10 seconds since the script started
		if time.Since(now).Seconds() >= 10 {
			break
		}
	}
}

func removeScript(sc script.Script, running *[]script.Script) {
	var newRunning []script.Script
	for _, s := range *running {
		if s.GetId() != sc.GetId() {
			newRunning = append(newRunning, s)
		}
	}
	*running = newRunning
}

// remove script from path and log from path
func removeScriptFiles(sc script.Script, db *db.Database, acc account.Account) {
	// get a list of all script ids
	ids := db.GetScriptIdsFromUid(acc.GetUser().GetId())
	var remove bool = true

	for _, id := range ids {
		if id == sc.GetId() {
			remove = false
			break
		}
	}

	if remove {
		// remove the script from the path
		os.Remove("userdata/" + acc.GetUser().GetUsername() + "/scripts/" + sc.GetName() + ".lua")
		// remove the log from the path
		os.Remove("userdata/" + acc.GetUser().GetUsername() + "/logs/" + sc.GetName() + ".log")
	}
}
