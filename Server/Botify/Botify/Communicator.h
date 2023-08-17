#pragma once

#pragma comment (lib, "ws2_32.lib")

#include <WinSock2.h>
#include <Windows.h>

#include <string>
#include <thread>
#include <iostream>
#include <vector>
#include <map>

#include "IRequestHandler.h"
#include "RequestHandlerFactory.h"
#include "ICryptoAlgorithm.h"

class Communicator
{
private:
	SOCKET _socket;
	std::map<SOCKET, IRequestHandler*> _clients;
	RequestHandlerFactory& m_handlers_factory;
	
	std::string _ip;
	int _port;

private:
	// does bindAndListen
	void create_socket();
	void bind_socket();
	
	void handle_new_client(SOCKET socket);
public:
	Communicator(const std::string& ip, const int port, RequestHandlerFactory& handlers_factory);
	~Communicator();

	void start_handle_requests();
	std::string get_logged_in();

private:

	void sign_user_out(SOCKET& socket);

	static void send_message(SOCKET socket, const std::vector<uint8_t> bytes);
	static std::vector<uint8_t> receive_message(SOCKET socket, const int numBytes);
};

