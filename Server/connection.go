package main

// this file will handle new connections to the server

import (
	"botify/main/db"
	"botify/main/logger"
	"botify/main/packer"
	um "botify/main/userman"
	"encoding/binary"
	"fmt"
	"net"
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
	var user um.User

	for {
		var header int16
		err := binary.Read(conn, binary.BigEndian, &header)
		if err != nil {
			logger.Error(err)
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
			HandleLogin(request, conn, &user, database)
		case SIGNUP:
			HandleSignup(request, conn, &user, database)
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

func HandleLogin(request packer.Request, conn net.Conn, user *um.User, database *db.Database) {
	// get the username and the password
	var username = request.Data["username"]
	var password = request.Data["password"]

	strUsername := fmt.Sprint(username)
	strPassword := fmt.Sprint(password)

	user.Logout()
	user.Login(strUsername, strPassword, *database)

	// centralize the response
	var response packer.Response
	response.Header = LOGIN
	response.Data = make(map[string]interface{})
	response.Data["confirmed"] = user.IsConfirmed()
	response.Data["id"] = user.GetId()
	response.Data["username"] = user.GetUsername()
	response.Data["email"] = user.GetEmail()
	response.Data["phone"] = user.GetPhone()

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)
}

func HandleSignup(request packer.Request, conn net.Conn, user *um.User, database *db.Database) {
	// get the username and the password
	var username = request.Data["username"]
	var password = request.Data["password"]
	var email = request.Data["email"]
	var phone = request.Data["phone"]

	strUsername := fmt.Sprint(username)
	strPassword := fmt.Sprint(password)
	strEmail := fmt.Sprint(email)
	strPhone := fmt.Sprint(phone)

	user.Logout()
	user.Signup(strUsername, strPassword, strEmail, strPhone, *database)

	// centralize the response
	var response packer.Response
	response.Header = SIGNUP
	response.Data = make(map[string]interface{})
	response.Data["confirmed"] = user.IsConfirmed()
	response.Data["id"] = user.GetId()
	response.Data["username"] = user.GetUsername()
	response.Data["email"] = user.GetEmail()
	response.Data["phone"] = user.GetPhone()

	// centralize the response
	data, _ := packer.Centralize(response)

	// send the response
	conn.Write(data)
}
