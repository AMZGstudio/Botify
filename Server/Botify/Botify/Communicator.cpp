#include "Communicator.h"
#include "LoginRequestHandler.h"

#include "JsonRequestPacketDeserializer.h"
#include "JsonResponsePacketSerializer.h"

#include "utils.h"
#include "SocketException.h"
#include "OTPCryptoAlgorithm.h"

#include <memory>

void Communicator::create_socket()
{
	// this server use TCP. that why SOCK_STREAM & IPPROTO_TCP
	// if the server use UDP we will use: SOCK_DGRAM & IPPROTO_UDP
	_socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

	if (_socket == INVALID_SOCKET)
		throw std::exception(__FUNCTION__ " - socket");
}

void Communicator::bind_socket()
{
	struct sockaddr_in sa = { 0 };

	sa.sin_port = htons(_port); // port that server will listen for
	sa.sin_family = AF_INET;		    // must be AF_INET
	sa.sin_addr.s_addr = INADDR_ANY;    // when there are few ip's for the machine. We will use always "INADDR_ANY"

	// Connects between the socket and the configuration (port and etc..)
	if (bind(_socket, (struct sockaddr*)&sa, sizeof(sa)) == SOCKET_ERROR)
		throw std::exception(__FUNCTION__ " - bind");
}

void Communicator::handle_new_client(SOCKET socket)
{
	Log::start_op("Created a thread to handle this socket", (int)socket);
	std::string client_username;

	// the encryption algorithm to use
	std::unique_ptr<ICryptoAlgorithm> encryption_algorithm = std::make_unique<OTPCryptoAlgorithm>();

	try
	{
		int protocol_code, message_size;
		
		_clients[socket] = this->m_handlers_factory.create_login_request();
		
		while (true)
		{
			// getting a request
			Log::start_op("Waiting for message...", (int)socket);
			protocol_code = 0, message_size = 0;

			// immediatly after receiving a message, decrypt it, for the rest of the program to use
			std::vector<uint8_t> header = encryption_algorithm->decrypt(receive_message(socket, 5));
			JsonRequestPacketDeserializer::deserialize_header(protocol_code, message_size, header);
			std::vector<uint8_t> data = encryption_algorithm->decrypt(receive_message(socket, message_size), 5);

			Log::end_op("Decrypted: [code:" + std::to_string(protocol_code)+", size:" + std::to_string(message_size)+", data: \'"+ std::string(data.begin(), data.end())+"\']", (int)socket);

			// handling that request
			RequestInfo ri { protocol_code, get_curr_time(), data, socket, encryption_algorithm };
			RequestResult rr = _clients.at(socket)->handle_request(ri);


			if (rr.newHandler != nullptr)
			{
				delete _clients[socket];
				_clients[socket] = rr.newHandler;
			}

			int response_code = 0, response_size = 0;
			JsonRequestPacketDeserializer::deserialize_header(response_code, response_size, rr.response);

			Log::msg("Sending Response: [code:"+std::to_string(response_code)+", data: \'"+ std::string(rr.response.begin()+5, rr.response.end())+"\']", (int)socket);
			
			// encrypt the message before sending it
			std::vector<uint8_t> encrypted_msg = encryption_algorithm->encrypt(rr.response);
			send_message(socket, encrypted_msg);
		}
	}

	catch (SocketException)
	{
		Log::error("Client forcefully disconnected, signing them out.", __FUNCTION__, (int)socket);
		sign_user_out(socket);
	}
	catch (std::exception& e)
	{
		Log::error(e.what(), __FUNCTION__, (int)socket);
	}
	
	closesocket(socket);
	Log::end_op("Closing thread for socket...", (int)socket);
}

Communicator::Communicator(const std::string& ip, const int port, RequestHandlerFactory& handlers_factory) : _ip(ip), _port(port), m_handlers_factory(handlers_factory)
{
	try
	{
		// this will inititalize the network
		WSADATA wsa_data = { };
		if (WSAStartup(MAKEWORD(2, 2), &wsa_data) != 0)
			throw std::exception("WSAStartup Failed");

		create_socket();
		bind_socket();
	}
	catch (std::exception& e)
	{
		MessageBoxA(NULL, e.what(), "Server Startup Error", MB_OK | MB_ICONERROR);
		exit(1);
	}
}

Communicator::~Communicator()
{
	try
	{
		// the only use of the destructor should be for freeing 
		// resources that was allocated in the constructor
		closesocket(_socket);

		// close the network
		WSACleanup();
	}
	catch (...) {}
}

void Communicator::start_handle_requests()
{
	// Start listening for incoming requests of clients
	if (::listen(_socket, SOMAXCONN) == SOCKET_ERROR)
		throw std::exception(__FUNCTION__ " - listen");

	Log::start_op(std::format("Listening on {}:{}", _ip, _port));

	while (true)
	{
		try
		{
			// this accepts the client and create a specific socket from server to this client
			// the process will not continue until a client connects to the server
			SOCKET client_socket = accept(_socket, NULL, NULL);
			if (client_socket == INVALID_SOCKET)
				throw SocketException("Accepted client socket is invalid!");

			try
			{
				// if a client connected, start a thread
				// note it is responsible for closing the socket when done.
				std::thread t(&Communicator::handle_new_client, this, client_socket);
				t.detach();
			}
			
			// if it crashed, then we will manually close the socket.
			catch (std::exception& e)
			{
				closesocket(client_socket);
				Log::error(e.what(), __FUNCTION__);
			}
		}
		catch (std::exception& e)
		{
			Log::error(e.what(), __FUNCTION__);
		}
	}

	Log::end_op("Finished listening.");
}

std::string Communicator::get_logged_in()
{
	return "TODO";
}

void Communicator::sign_user_out(SOCKET& socket)
{
	m_handlers_factory.get_login_manager().logout(socket);
	_clients.erase(socket);
}


void Communicator::send_message(SOCKET socket, const std::vector<uint8_t> bytes)
{	
	if (send(socket, std::string(bytes.begin(), bytes.end()).c_str(), (int)bytes.size(), 0) == INVALID_SOCKET)
		throw std::exception("Error while sending message to client");	
}

std::vector<uint8_t> Communicator::receive_message(SOCKET socket, const int numBytes)
{
	std::vector<uint8_t> bytes;

	if (numBytes == 0)
		return bytes;
		
	char* data = new char[numBytes + 1];
	
	if (recv(socket, data, numBytes, 0) == SOCKET_ERROR)
		throw SocketException((std::string("Error while recieving from socket: ") + std::to_string(socket)).c_str());
	
	for (int i = 0; i < numBytes; i++)
		bytes.push_back(data[i]);

	delete[] data;
	return bytes;
}