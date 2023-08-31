# in this script we will test the server
# send a request to the server and get the response and print it

import socket
import json
import enum


class Command(enum.Enum):
    LOGIN = enum.auto()
    SIGNUP = enum.auto()
    AUTHENTICATE = enum.auto()
    UPLOAD_SCRIPT = enum.auto()
    GET_SCRIPT = enum.auto()
    GET_SCRIPT_LIST = enum.auto()
    EDIT_SCRIPT = enum.auto()
    REMOVE_SCRIPT = enum.auto()
    ACTIVATE_SCRIPT = enum.auto()
    DISABLE_SCRIPT = enum.auto()
    TRIGGER = enum.auto()
    ACTION = enum.auto()

# so the messege looks like this:
# 1. the first 2 bytes are the header
# 2. the next 4 bytes are the length of the messege
# 3. the rest of the bytes are the messege which is a json object

# create a function that centralize a message by header and dict
def centralize_message(header, dict):
    # convert the header to bytes
    header = header.to_bytes(2, byteorder="big")
    # convert the dict to a json object
    json_object = json.dumps(dict)
    # convert the json object to bytes
    json_object_bytes = json_object.encode()
    # get the length of the json object
    json_object_length = len(json_object_bytes)
    # convert the length to bytes
    json_object_length_bytes = json_object_length.to_bytes(4, byteorder="big")
    # create the message
    message = header + json_object_length_bytes + json_object_bytes
    # return the message
    return message

# create a function that decetralize a message
def decentralize_message(message):
    # get the header
    header = message[0:2]
    # turn the header to an int
    header = int.from_bytes(header, byteorder="big")
    # get the length of the json object
    json_object_length = int.from_bytes(message[2:6], byteorder="big")
    # get the json object
    json_object = message[6:6 + json_object_length]
    # convert the json object to a dict
    try:
        dict = json.loads(json_object)
    except:
        dict = {}
    # return the header and the dict
    return header, dict

# send a request to the server and get the response and print it

# create a socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# connect to the server
s.connect(("127.0.0.1", 8080))

# we want to login
request = centralize_message(Command.LOGIN.value-1, {"username": "lavi", "password": "1111"})
print(request)
# send the request
s.send(request)
# get the response
response = s.recv(1024)
# decentralize the response
header, dict = decentralize_message(response)
# print the response
print(dict)

# create a menu
while True:
    # 1 - login
    # 2 - signup
    # else - exit
    choice = input("1 - login\n2 - signup\n3 - upload script\n4 - remove script\nelse - exit\n")
    if choice == "1":
        # we want to login
        username = input("username: ")
        password = input("password: ")
        request = centralize_message(Command.LOGIN.value-1, {"username": username, "password": password})
        # send the request
        print(request)
        s.send(request)
        # get the response
        response = s.recv(1024)
        # decentralize the response
        header, dict = decentralize_message(response)
        # print the response
        print(dict)
    elif choice == "2":
        # we want to signup
        username = input("username: ")
        password = input("password: ")
        request = centralize_message(Command.SIGNUP.value-1, {"username": username, "password": password})
        # send the request
        print(request)
        s.send(request)
        # get the response
        response = s.recv(1024)
        # decentralize the response
        header, dict = decentralize_message(response)
        # print the response
        print(dict)
    elif choice == "3":
        # we want to upload a script
        script_name = input("script name: ")
        script = 'print("hello mister")'
        request = centralize_message(Command.UPLOAD_SCRIPT.value-1, {"scriptName": script_name, "script": script})
        # send the request
        print(request)
        s.send(request)
        # get the response
        response = s.recv(1024)
        # decentralize the response
        header, dict = decentralize_message(response)
        # print the response
        print(dict)
    elif choice == "4":
        # we want to get a script
        script_name = input("script name: ")
        request = centralize_message(Command.REMOVE_SCRIPT.value-1, {"scriptName": script_name})
        # send the request
        print(request)
        s.send(request)
        # get the response
        response = s.recv(1024)
        # decentralize the response
        header, dict = decentralize_message(response)
        # print the response
        print(dict)
    else:
        # we want to exit
        break