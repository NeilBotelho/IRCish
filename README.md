# IRCish

This project aims to implement a chat server in the style of irc, where clients connect to rooms and recieve messages sent to those rooms only while connected.

The project will be attempted using golang and websockets.

## Current progress
Base of the server taken from my [basic golang server](https://github.com/NeilBotelho/basic-golang-server/). Resilliency added to server by setting timeouts for reads and writes. 

~~The plan is to use gorilla/websocket for websockets.~~

>Update 20/07/2020

Basic websocket connection implemented using gorilla/websocket