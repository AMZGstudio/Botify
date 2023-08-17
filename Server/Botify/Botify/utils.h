#pragma once

#include <string>
#include <iostream>
#include <ctime>
#include <sstream>
#include <iomanip>

#include "ICryptoAlgorithm.h"
#include "RequestAndResponseStructs.h"

#include "nlohmann/json.hpp"
using Json = nlohmann::json;

enum Code { Error, Success, Login, Signup, Logout, GetRooms, getPlayersInRoom, JoinRoom, CreateRoom, HighScore, PersonalStats, CloseRoom, LeaveRoom, StartGame, RoomState, SubmitAnswer, LeaveGame, GameResult, GetQuestion, AddQuestion };

enum QuestionStatus { MoreQuestions, NoMoreQuestions};
enum GameStatus { GameNotOver, GameOver};

#define COLOR_RED "31"

class Log
{
public:

	static std::string add_color(std::string text, std::string color_code);

	static void error(std::string err, const char* function_name, int sock_id = 0);
	static void start_op(std::string msg, int sock_id = 0);
	static void end_op(std::string msg, int sock_id = 0);
	static void msg(std::string msg, int sock_id = 0);
};

std::string get_curr_time();