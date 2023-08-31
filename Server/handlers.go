package main

// this file will handle new connections to the server

import (
	acc "botify/main/account"
	"botify/main/db"
	"botify/main/logger"
	"botify/main/packer"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	LOGIN = iota
	SIGNUP
	AUTHENTICATE
	UPLOAD_SCRIPT
	GET_SCRIPT
	GET_SCRIPT_LIST
	EDIT_SCRIPT
	REMOVE_SCRIPT
	ACTIVATE_SCRIPT
	DISABLE_SCRIPT
	TRIGGER
	ACTION
)

func HandleConnection(conn net.Conn, database *db.Database) {
	// create a user
	var account acc.Account

	for {
		var header int16
		err := binary.Read(conn, binary.BigEndian, &header)
		if err != nil {
			logger.Info("connection closed by " + conn.RemoteAddr().String())
			conn.Close()
			return
		}

		var length int32
		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			logger.Error(err)
			return
		}

		buf := make([]byte, length)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		headerBuf := make([]byte, 2)
		lengthBuf := make([]byte, 4)

		binary.BigEndian.PutUint16(headerBuf, uint16(header))
		binary.BigEndian.PutUint32(lengthBuf, uint32(length))

		fullBuf := append(headerBuf, lengthBuf...)
		fullBuf = append(fullBuf, buf[:n]...)

		request, _ := packer.Decentralize(fullBuf)

		switch request.Header {
		case LOGIN:
			HandleLogin(request, conn, &account, database)

		case SIGNUP:
			HandleSignup(request, conn, &account, database)

		case UPLOAD_SCRIPT:
			if account.GetUser().IsConfirmed() && account.GetUser().GetConnectionType() == "control" {
				HandleUploadScript(request, conn, &account, database)
			} else {
				unauthorizedAccess(conn)
			}

		// edit a script
		case EDIT_SCRIPT:
			if account.GetUser().IsConfirmed() && account.GetUser().GetConnectionType() == "control" {
				HandleUploadScript(request, conn, &account, database)
			} else {
				unauthorizedAccess(conn)
			}

		case REMOVE_SCRIPT:
			if account.GetUser().IsConfirmed() && account.GetUser().GetConnectionType() == "control" {
				HandleRemoveScript(request, conn, &account, database)
			} else {
				unauthorizedAccess(conn)
			}

		default:
			// send back "hello world"
			var response packer.Response
			response.Header = 0
			response.Data = make(map[string]interface{})
			response.Data["hello"] = "world"

			// centralize the response
			data, _ := packer.Centralize(response)

			// send the response
			conn.Write(data)
		}
	}
}

func HandleLogin(request packer.Request, conn net.Conn, account *acc.Account, database *db.Database) {
	// get the username and the password
	var username = request.Data["username"]
	var password = request.Data["password"]
	var connectionType = request.Data["connectionType"]

	strUsername := fmt.Sprint(username)
	strPassword := fmt.Sprint(password)
	strConnectionType := fmt.Sprint(connectionType)

	// check if we have anything in the connection type
	if connectionType == nil {
		strConnectionType = "control"
	}

	logger.Info("handling login of user: " + strUsername + " with password: " + strPassword + " and connection type: " + strConnectionType)

	// create a user
	var user acc.User

	user.Logout()
	user.Login(strUsername, strPassword, strConnectionType, *database)
	*account = acc.CreateAccount(user, *database)

	// centralize the response
	var response packer.Response
	response.Header = LOGIN
	response.Data = make(map[string]interface{})
	response.Data["confirmed"] = user.IsConfirmed()
	response.Data["id"] = account.GetId()
	response.Data["scripts"] = account.GetScriptNames()
	response.Data["username"] = user.GetUsername()
	response.Data["connectionType"] = user.GetConnectionType()

	// if login was successful, start chrome
	if user.IsConfirmed() && user.GetConnectionType() == "service" {
		go startChrome(conn, account)
	}

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)

}

func HandleSignup(request packer.Request, conn net.Conn, account *acc.Account, database *db.Database) {
	// get the username and the password
	var username = request.Data["username"]
	var password = request.Data["password"]
	var email = request.Data["email"]
	var phone = request.Data["phone"]

	strUsername := fmt.Sprint(username)
	strPassword := fmt.Sprint(password)
	strEmail := fmt.Sprint(email)
	strPhone := fmt.Sprint(phone)
	strConnectionType := "control"

	logger.Info("handling signup of user: " + strUsername + " with password: " + strPassword)

	// create a user
	var user acc.User

	user.Logout()
	user.Signup(strUsername, strPassword, strEmail, strPhone, strConnectionType, *database)
	*account = acc.CreateAccount(user, *database)

	// centralize the response
	var response packer.Response
	response.Header = SIGNUP
	response.Data = make(map[string]interface{})
	response.Data["confirmed"] = user.IsConfirmed()
	response.Data["id"] = account.GetId()
	response.Data["scripts"] = account.GetScriptNames()
	response.Data["username"] = user.GetUsername()
	response.Data["connectionType"] = user.GetConnectionType()

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)

}

func HandleUploadScript(request packer.Request, conn net.Conn, account *acc.Account, database *db.Database) {
	// get the script name and the script
	var scriptName = request.Data["scriptName"]
	var script = request.Data["script"]

	strScriptName := fmt.Sprint(scriptName)
	strScript := fmt.Sprint(script)

	logger.Info("handling upload of script: " + strScriptName + " by user: " + account.GetUser().GetUsername())

	account.GetUser().CreateFolders()
	var path = "userdata/" + account.GetUser().GetUsername() + "/scripts/" + strScriptName + ".lua"

	// open file and write the script
	file, err := os.Create(path)
	if err != nil {
		logger.Error(err)
		return
	}
	// write the script to the file
	file.WriteString(strScript)

	defer file.Close()

	// now add the script to the database
	account.AddScript(strScriptName, "", path, *database)

	// centralize the response
	var response packer.Response
	response.Header = UPLOAD_SCRIPT
	response.Data = make(map[string]interface{})
	response.Data["confirmed"] = true

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)

}

func unauthorizedAccess(conn net.Conn) {
	// send back "unauthorized access"
	var response packer.Response
	response.Header = 0
	response.Data = make(map[string]interface{})
	response.Data["error"] = "unauthorized access"

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)
}

// remove a script
func HandleRemoveScript(request packer.Request, conn net.Conn, account *acc.Account, database *db.Database) {
	// get the script id
	var scriptName = request.Data["scriptName"]
	strScriptName := fmt.Sprint(scriptName)
	scriptId := database.GetIdFromName(strScriptName, "scripts")

	if scriptId == -1 {
		logger.Warn("script: " + strScriptName + " was not found")

		var response packer.Response
		response.Header = 0
		response.Data = make(map[string]interface{})
		response.Data["error"] = "script not found"

		data, _ := packer.Centralize(response)
		conn.Write(data)

		return
	}

	logger.Info("handling remove of script: " + strScriptName + " by user: " + account.GetUser().GetUsername())

	// remove the script from the database
	account.DeleteScript(scriptId, *database)

	// centralize the response
	var response packer.Response
	response.Header = REMOVE_SCRIPT
	response.Data = make(map[string]interface{})
	response.Data["confirmed"] = true

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)
}

// a function for the poc that will get a connection to the server and send to it: {"command":"start chrome"}
func startChrome(conn net.Conn, acc *acc.Account) {
	for {
		time.Sleep(5 * time.Second)
		// if connection type is service, send a command {"command":"start chrome"}
		if acc.GetUser().GetConnectionType() != "service" {
			return
		}

		// send a command to the server
		var response packer.Response
		response.Header = ACTION
		response.Data = make(map[string]interface{})
		response.Data["command"] = "start chrome"

		// centralize the response
		data, _ := packer.Centralize(response)

		// send the response
		conn.Write(data)
	}
}
