# IRCish

## What is this
This project aims to implement a chat server in the style of irc, where clients connect to rooms and recieve messages sent to those rooms only while connected.

The project will be attempted using golang and websockets.

Current Progress can be found in PROGRESS.md

## How can I try it?
<!--At the time of writing this(Jul 21 2020 12:31AM) there is a live version of with all the current features at https://ircish.herokuapp.com/-->

The project is still in development so the instructions below may not be up to date(I'll try my best ;) ) but at the time of writing this, to run the project:

1. Clone this repository
2. Run ``` go run *.go``` in the top level of this repository
3. Open the index.html file in your browser 

and voila you have the IRCish server running locally and every new tab or window that you open to the index.html file is a new client to the server and you can send messages between them
