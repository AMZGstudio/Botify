#include "SocketException.h"

SocketException::SocketException(const std::string err_message) : message(err_message)
{
}

std::string SocketException::what()
{
	return message;
}
