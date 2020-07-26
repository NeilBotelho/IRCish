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

- ~Break the wsHandler function into smaller function(its getting a bit long and unwieldly)~
- ~implement identify logic~
- implement ping logic and enable read timeouts on sockets
- Update client side logic and UI to handle json and multiple rooms(will probably take me the longest)
  
#### 26 Jul 2020
- Client UI completed. Took a while, not very good at frontend. 
- Client UI made responsive
- Rooms implemented in client UI
- Identify logic implemented server side

**TODO**
- Complete client side logic for:
	- Sending and recieving messages
	- Changing identity
	- leaving and joining rooms
	- switching rooms (will take a while)
