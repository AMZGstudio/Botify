#pragma once

#include "Communicator.h"
#include "RequestHandlerFactory.h"

class Server
{
private:
	// keep this the member order!
	std::shared_ptr<IDatabase> m_database;
	RequestHandlerFactory m_handlers_factory;
	Communicator _communicatorThread;

public:
	Server();
	void print_logged_in();
	void run();
};

