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
			for _, s := range a.GetScripts() {
				// try to find the script in the running list
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
					go runScript(s, a, &running)
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
	*accPtr = acc
}

func runScript(sc script.Script, acc account.Account, running *[]script.Script) {
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

		if err := L.DoFile("userdata/" + username + "/scripts/" + sc.GetName() + ".lua"); err != nil {
			logger.Error(err)
			return
		}

		time.Sleep(500 * time.Millisecond)

		// check if it has been 1 minute
		if time.Since(now).Minutes() >= 1 {
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
