# IRCish

## What is this
This project implements a chat server in the style of irc, where clients connect to rooms and recieve messages sent to those rooms only while connected. It is implemented using golang and websockets.
Current Progress can be found in PROGRESS.md

## Commands available
- __/join__ _roomName_ 
to join the roomwith name __roomName__. Roon names must contain only lowercase letters, numbers and underscores. Room names must be between 2 and 10 characters long
- __/identify__ _username_
to change how you are identified. Usernames can contain any case letters, numbers and must be between 2 and 10 characters long
- __/clear__
clear all messages in current room
- /help
prints a list of available commands


## How can I try it?
At the time of writing this(Jul 29 2020 8:55PM) there is a live version of with all the current features at https://ircish.herokuapp.com/



## Running Locally
1. Clone this repository
2. Run ``` go build && ./server``` in the top level of this repository
3. Open the localhost:8000 file in your browser 

and voila you have the IRCish server running locally and every new tab or window that you open to the index.html file is a new client to the server and you can send messages between them
