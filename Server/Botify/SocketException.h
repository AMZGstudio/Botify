#pragma once
#include <iostream>

class SocketException : public std::exception
{
private:
	std::string message;

public:
	SocketException(const std::string err_message);

	std::string what();
};

