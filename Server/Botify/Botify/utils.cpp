#include "utils.h"

std::string get_socket_string(int sock_id)
{
	if (sock_id == 0)
		return "";

	return "s[" + std::to_string(sock_id) + "]-";
}
std::string Log::add_color(std::string text, std::string color_code)
{
	return "\033[" + color_code + 'm' + text + "\033[0m";
}

void Log::error(std::string err, const char* function_name, int sock_id)
{
	std::string s = get_socket_string(sock_id) + add_color("Error  : " + err + " in function: " + function_name, COLOR_RED) + '\n';
	std::cerr << s;
}

void Log::start_op(std::string msg, int sock_id)
{
	std::string s = get_socket_string(sock_id) + "Running: " + msg + '\n';
	std::cout << s;
}

void Log::end_op(std::string msg, int sock_id)
{
	std::string s = get_socket_string(sock_id) + "Result : " + msg + '\n';
	std::cout << s;
}

void Log::msg(std::string msg, int sock_id)
{
	std::string s = get_socket_string(sock_id) + msg + '\n';
	std::cout << s;
}


inline std::tm localtime_xp(std::time_t timer)
{
	std::tm bt {};
#if defined(__unix__)
	localtime_r(&timer, &bt);
#elif defined(_MSC_VER)
	localtime_s(&bt, &timer);
#else
	static std::mutex mtx;
	std::lock_guard<std::mutex> lock(mtx);
	bt = *std::localtime(&timer);
#endif
	return bt;
}

std::string get_curr_time()
{
	auto t = std::time(nullptr);
	auto tm = localtime_xp(t);

	std::ostringstream oss;
	oss << std::put_time(&tm, "%d-%m-%Y %H-%M-%S");
	return oss.str();
}
