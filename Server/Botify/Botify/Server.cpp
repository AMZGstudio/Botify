#include "Server.h"
#include "SqliteDatabase.h"

Server::Server() : m_database(std::make_shared<SqliteDatabase>("sql_database.sqlite")),  m_handlers_factory(m_database), _communicatorThread("127.0.0.1", 9090, m_handlers_factory)
{
}

void Server::print_logged_in()
{
	std::cout << _communicatorThread.get_logged_in() << std::endl;
}

void Server::run()
{
	// create a thread that runs the communicators listen function, and detach it
	std::thread connect(&Communicator::start_handle_requests, std::ref(_communicatorThread));
	connect.detach();
}
