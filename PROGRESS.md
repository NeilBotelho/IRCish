# Progress So Far

#### 19 Jul 2020
Base of the server taken from my [basic golang server](https://github.com/NeilBotelho/basic-golang-server/). Resilliency added to server by setting timeouts for reads and writes. 

The plan is to use gorilla/websocket for websockets.

#### 20 Jul 2020
- Basic websocket connection implemented using gorilla/websocket on server side
- Basic functioning front end built(message list, input box)
- Websocket functionality implemented on client side. 
- Server and client can now send messages to each other. But currently one client cannot send a message to another 
- Implemented communication between clients. Rooms/channels not yet implemented. Only a single room exists

**TODO**

- ~~find how to identify disconnected clients~~
- ~~fix memory leak of stale clients~~

#### 21 Jul 2020
- Rewrite of broadcast and message system complete(see DesignDoc.md for details on implementation)
- JSON communication with client implemented on server side(client side remaining) 
- Broke client.go into room.go(Room, message and braodcast definitions) and client.go(client handlers and struct definitions)
- Rooms functionality implemented
- Users can now join and leave rooms as they wish
- Users are now identified by a randomly generated number rather than their public IP

**TODO**

- ~~Break the wsHandler function into smaller function(its getting a bit long and unwieldly)~~
- ~~implement identify logic~~
- ~~implement ping logic and enable read timeouts on sockets~~
- ~~Update client side logic and UI to handle json and multiple rooms(will probably take me the longest)~~
  
#### 26 Jul 2020
- Client UI completed. Took a while, not very good at frontend. 
- Client UI made responsive
- Rooms implemented in client UI
- Identify logic implemented server side

**TODO**
- ~~Complete client side logic for:~~
	- ~~Sending and recieving messages~~
	- ~~Changing identity~~
	- ~~leaving and joining rooms~~
	- ~~switching rooms (will take a while)~~

### 29 Jul 2020
(I was sick for the past week and half hence the slow progress)
- Ping logic implemented(server and client) and read timeouts enabled
	With ping logic implemented, on mobile even if the user moves to another app(without purging browser cache or closing the tab) the connection remains alive and there is no need to reload.
- Client side logic updated to handle multiple rooms
- Client side UI updated to handle multiple rooms
- Client side logic updated to handle /identify and /clear commands
- Client side UI refreshed to look better
- Client side UI updated to notify when a room has new messages

----
## Note
At this point this project is essentially complete. All initial goals have been achieved. But along the way I thought of a few improvements that could be made. I may or may not implement these moving forward, they will not be a priority.

**Possible Improvements**
1. Enable reconnect when connection is lost

	Instead of a user having to reload when for some reason he/she disconnects from the server(eg. internet loss) the user can issue a ```/reconnect``` command and be reconnected, keeping all their past messages. This requires the client to tell the server the users identity as well as the rooms the client was connected to previously without the server re-announcing its entry. This would require one or more new entries in Msg and a new opcode  
