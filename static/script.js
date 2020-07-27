identified=false//check if we've identified with server only then allow sendmessage
roomRegex=RegExp('[a-z\-]{3,10}')
//Front end implementation
var roomsList=['general']
var messageFeed=document.getElementById('message-feed')
var currentRoom=document.getElementById('general-messages')
currentRoom.setAttribute("style","display:block")
var inputField=document.getElementById('input-field')
var help="AVAILABLE COMMANDS:\n"+
"/join channelName : to join a channel(name must contain only lowercase letters and underscores and must begin with a letter)\n"

function updateMessages(newline,room="general",classes=[]){
	messages=document.getElementById(room+"-messages")
	var newcontent = document.createElement('p');
	newcontent.innerText = newline+"\n";
	while (newcontent.firstChild) {
		   messages.appendChild(newcontent.firstChild);
	}
	for(var c in classes){
		newcontent.classList.add(c)
	}
	messageFeed.scrollTop = messageFeed.scrollHeight;
}

function createRoom(roomName){
	// Create containing div
	var newRoom = document.createElement('div');
	newRoom.classList.add("message-display")
	newRoom.id=roomName+"-messages"
	messageFeed.appendChild(newRoom);
	// Add paragraph tag to contain messages
	var roomBody= document.createElement('p');
	while (roomBody.firstChild) {
   		newRoom.appendChild(roomBody.firstChild);
	}

	// Create room in sidebar
	var roomButton=document.createElement('p')
	roomButton.classList.add("room-name")
	roomButton.id=roomName+"-room"
	roomButton.innerText="#"+roomName
	//append new button to roomList
	roomtray=document.getElementById("room-list")
	roomtray.appendChild(roomButton)
	roomsList.push(roomName)
}

function roomSwitch(roomName){
	currentRoom.setAttribute("style","display:none")
	currentRoom=document.getElementById(roomName+"-messages")
	currentRoom.setAttribute("style","display:block")

}

function sendMessage(e){
	// send message to server
	let msg=e.target.value.trim()
	if(msg.substring(0,1)=="/"){
		if(msg.substring(0,5)=="/join"){
			tokens=msg.split(" ")
			if (tokens.length==2 && roomRegex.test(tokens[1])){
				// Is a legitamate join
				if(roomsList.indexOf(tokens[1])==-1){
					ws.send(JSON.stringify({"opCode":1,"room":tokens[1]}))
					console.log("MEssage"+tokens[1])
					inputField.value=""
					return
				}
			}
		}
		if(msg=="/help"){
			updateMessages(help,['system-notification'])
			inputField.value=""
			return
		}
	}
	
	if(ws.readyState != WebSocket.OPEN){
		updateMessages("Not connected to server")
	}
	else{
		updateMessages(e.target.value)
		console.log(e.target.value)
		// console.log(e.target.value)
		//send message over socket
		// ws.send(JSON.stringify({"opCode":0,"content":e.target.value}))
		// console.log()
	}
	inputField.value=""
}

// Register input field onchange function
inputField.onchange=sendMessage

// Websockets implementation
// Check if websockets are supported 
if(!(window.WebSocket)){
	alert("Websockets not supported in this browser")
}

// Create websocket and register handlers
ws = new WebSocket("ws://localhost:8000/ws");
ws.onerror=function(event){
	displayError()
}
ws.onopen=function(event){
	console.log("connected to server")
}
ws.onmessage=function(event){
	// console.log(event.data)
	reply=JSON.parse(event.data)
	switch(reply.opcode){
		case 0:
		updateMessages(reply.from+": "+reply.content)
		break
		case 1:

	} 
	console.log(event.data)
	// updateMessages(event.data,['system-notification'])
}
function displayError(){
	alert("Error connecting to server")
}
createRoom("yoyo")