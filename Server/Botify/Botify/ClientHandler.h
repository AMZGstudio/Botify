#pragma once

#include "Server.h"

class ClientHandler
{
private:
	SOCKET _socket;
	std::string _name;

private:

public:
	ClientHandler(SOCKET socket);
	void sendMessage(const std::string str);
	std::string receiveMessage(const int numBytes);
};

void startClientHandler(SOCKET socket);
