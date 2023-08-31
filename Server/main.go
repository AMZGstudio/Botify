package main

import (
	"botify/main/db"
	"botify/main/logger"
	"net"
)

func main() {
	// create a new database
	var database db.Database
	database.SetPath("db.sqlite")
	go handleAccounts(&database)
	// wait for a new connection
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error(err)
		return
	}
	for {
		// accept the connection
		conn, err := listen.Accept()
		if err != nil {
			logger.Error(err)
			continue
		}

		logger.Info("New connection from " + conn.RemoteAddr().String())
		go HandleConnection(conn, &database)
	}
}
