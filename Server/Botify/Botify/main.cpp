#include "Server.h"

int main()
{
	Server server;
	server.run();

	Log::msg("Type 'EXIT' to close the server\n");

	while (1)
	{
		std::string command;
		std::cin >> command;

		if (command == "EXIT")
			return 0;

		if (command == "QUERY")
			server.print_logged_in();

		if (command == "CLS")
			system("cls");

		// if (command == "SEND")

	}
}