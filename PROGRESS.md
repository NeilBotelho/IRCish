# Progress So Far

#### 19 Jul 2020
Base of the server taken from my [basic golang server](https://github.com/NeilBotelho/basic-golang-server/). Resilliency added to server by setting timeouts for reads and writes. 

The plan is to use gorilla/websocket for websockets.

#### 20 Jul 2020
- Basic websocket connection implemented using gorilla/websocket on server side
- Basic functioning front end built(message list, input box)
- Websocket functionality implemented on client side. 
- Server and client can now send messages to each other. But currently one client cannot send a message to another 

