package main

import (
	"botify/main/db"
	"botify/main/logger"
	"net"
)

/*User	Represents a user of the service
Script	Represents a script that can be run on the server
Device	Represents a device that can be used to get data
Account	Represents an online account that can be used to get data
Action	Represents an action that can be performed
Database	Represents the database used to store user data
Request	Represents a request made by a user
Response	Represents a response sent back to the user
Resource	Represents a resource that can be accessed via the API
Endpoint	Represents a URL endpoint for accessing a resource
Server	Represents the server that runs the service
Client	Represents a client that interacts with the service
Service	Represents the service itself
UserManager	Manages user accounts and authentication
StateManager	Manages the state of the service
ScriptManager	Manages scripts that can be run on the server
DeviceManager	Manages devices that can be used to get data
AccountManager	Manages online accounts that can be used to get data
ActionManager	Manages actions that can be performed*/

func main() {
	// create a new database
	var database db.Database
	database.SetPath("db.sqlite")
	go handleAccounts(&database)
	// wait for a new connection
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		// accept the connection
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		logger.Info("New connection from " + conn.RemoteAddr().String())
		go HandleConnection(conn, &database)
	}
}
