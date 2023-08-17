#include "ClientHandler.h"

/*

This is the client handler. it has a socket, and a name. you can call login, and loop on it.

*/

ClientHandler::ClientHandler(SOCKET socket) : _socket(socket)
{
}

void ClientHandler::sendMessage(const std::string str)
{
	if (send(_socket, str.c_str(), str.length(), 0) == INVALID_SOCKET)
	{
		throw std::exception("Error while sending message to client");
	}
}

std::string ClientHandler::receiveMessage(const int numBytes)
{
	if (numBytes == 0)
	{
		return "";
	}

	char* data = new char[numBytes + 1];
	int res = recv(_socket, data, numBytes, 0);
	if (res == INVALID_SOCKET)
	{
		std::string s = "Error while recieving from socket: ";
		s += std::to_string(_socket);
		throw std::exception(s.c_str());
	}
	data[numBytes] = 0;
	std::string received(data);
	delete[] data;
	return received;
}

void startClientHandler(SOCKET socket)
{
	ClientHandler ch(socket);

	try
	{
		std::cout << "Sending: Hello" << std::endl;
		ch.sendMessage("Hello");
		std::string str = ch.receiveMessage(5);
		std::cout << "Recieved: " + str << std::endl;
	}
	catch (...)
	{
		// Closing the socket (in the level of the TCP protocol)
		closesocket(socket);
	}
}