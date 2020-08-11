//FRONT END IMPLEMENTATION

validRoomName=RegExp('[a-z0-9\-]{2,10}')
validUsername=RegExp('[a-zA-z0-9\-]{2,10}')
var roomsList=[] //List of joined rooms
//Outer div for all  rooms messages
var messageFeed=document.getElementById('message-feed')//used for scrolling to the top
var inputField=document.getElementById('input-field')//used for clearing input field
var currentMessages
var currentRoom
var help="AVAILABLE COMMANDS:\n\n"+
"/join channelName : to join a room. Room names must contain only lowercase letters, numbers and underscores. Room names must be between 2 and 10 characters long)\n\n"+
"/identify username : to change how you are identified to 'username'. Usernames can contain any case letters, numbers and underscores. Usernames must be between  2 and 10 characters long. Usernames are not unique\n\n"+
"/leave : leaves the current room. Removes all messages\n"+
"/clear : clears all messages from current room"

function getRoomFromId(roomId){return roomId.split("-")[0]}

function createRoom(roomName){
	if(roomsList.indexOf(roomName)!=-1){
		return
	}
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
	roomButton.addEventListener("click",(e)=>{roomSwitch(e.target.id)})

	//append new button to roomList
	roomtray=document.getElementById("room-list")
	roomtray.appendChild(roomButton)
	roomsList.push(roomName)
}

function deleteRoom(roomName){
	roomsList.splice(roomsList.indexOf(roomName),1)		
	if(roomsList.indexOf('general')==-1){
		if(roomsList.length!=0){
			roomSwitch(roomsList[0])
		}
	}
	else{
		roomSwitch('general')
	}
	roomButton=document.getElementById(roomName+"-room")
	roomBody=document.getElementById(roomName+"-messages")
	roomtray=document.getElementById("room-list")

	roomtray.removeChild(roomButton)
	messageFeed.removeChild(roomBody)

}
function roomSwitch(roomId){
	roomName=getRoomFromId(roomId)//from roomId to roomName
	if(currentMessages!=undefined){
	currentMessages.setAttribute("style","display:none")
	currentRoom.classList.remove("active-room")
	}
	currentMessages=document.getElementById(roomName+"-messages")
	currentMessages.setAttribute("style","display:block")
	
	currentRoom=document.getElementById(roomName+"-room")
	currentRoom.classList.remove("new-messages")
	currentRoom.classList.add("active-room")
}


function updateMessages(newline,roomName=null,classes=[]){
	// create new message in the room's messages
	currRoomName=getRoomFromId(currentMessages.id)
	if(roomName==null){
		messages=document.getElementById(currRoomName+"-messages")
	}
	else{
		messages=document.getElementById(roomName+"-messages")
		if (currRoomName!=roomName)
		document.getElementById(roomName+"-room").classList.add("new-messages")
	}
	var newcontent = document.createElement('p');
	newcontent.innerText = newline+"\n";
	//Push to end
	while (newcontent.firstChild) {
		   messages.appendChild(newcontent.firstChild);
	}
	// add classes
	classes.forEach((className)=>{
		newcontent.classList.add(className)
	})

	// Scroll to bottom
	messageFeed.scrollTop = messageFeed.scrollHeight;
}

function sendMessage(e){
	// send message to server
	let msg=e.target.value.trim()
	let room=getRoomFromId(currentRoom.id)
	if(msg.substring(0,5)=="/join"){
		tokens=msg.split(" ")
		if (tokens.length==2 && validRoomName.test(tokens[1])){
			// Is a legitamate roomName
			ws.send(JSON.stringify({"opcode":1,"room":tokens[1]}))
			createRoom(tokens[1])
			roomSwitch(tokens[1]+"-room")
		}
		else{
			updateMessages("Invalid room name")
		}
		inputField.value=""
		return
	}
	if(msg=="/help"){
		updateMessages(help,roomName=null,classes=['system-notification'])
		inputField.value=""
		return
	}
	
	if(msg.substring(0,9)=="/identify"){
		tokens=msg.split(" ")
		if(tokens.length==2 && validUsername.test(tokens[1]) && tokens[1]!="System"){
			ws.send(JSON.stringify({"opcode":3,"content":tokens[1]}))
		}
		else{
			updateMessages("Invalid username")
		}
		inputField.value=""
		return

	}
	if(msg=="/clear"){
		currentMessages.innerText=""
		inputField.value=""
		return
	}
	if(msg=="/leave"){
		currRoomName=getRoomFromId(currentRoom.id)
		ws.send(JSON.stringify({"opcode":2,"room":getRoomFromId(currentRoom.id)}))
		deleteRoom(currRoomName)
		inputField.value=""
		return
	}

	if(msg.substring(0,1)=="/"){
		updateMessages("System: Messages cannot start with a '/' only commands can. Use /help for a list of available commands")
		inputField.value=""
		return
	}
	
	// send message over socket
	ws.send(JSON.stringify({"opcode":0,"content":e.target.value,"room":room}))
	inputField.value=""
}



// Websockets implementation
// Check if websockets are supported 
if(!(window.WebSocket)){
	alert("Websockets not supported in this browser")
}

// Create websocket and register handlers
// ws = new WebSocket("ws://localhost:8000/ircish"); // For local development
ws = new WebSocket("wss://ircish.herokuapp.com/ircish"); // For local development

ws.onerror=function(event){
	displayError()
}
ws.onopen=function(event){
	console.log("connected to server")
}
ws.onclose=function(event){
	alert("Connection to server closed")
}
ws.onmessage=function(event){
	// console.log(event.data)
	reply=JSON.parse(event.data)
	// console.log(reply)
	switch(reply.opcode){
		case 0:
		content=reply.from+": "+reply.content
		updateMessages(content,reply.room)
		break
		case 4:
		ws.send(JSON.stringify({"opcode":4}))
		break
		case 6:
		if(reply.room!=null){
			updateMessages("System: "+reply.content,roomName=reply.room,classes=['system-notification'])
		}
		else{
			updateMessages("System: "+reply.content,roomName=null,classes=['system-notification'])
		}
		break
		default:
		console.log("No handle case",reply)
	} 
}
function displayError(){
	alert("Error connecting to server")
}


//  MAIN 

// Register input field onchange function
inputField.onchange=sendMessage
//Create general room and set it to current Room
createRoom("general")
currentMessages=document.getElementById('general-messages')
currentRoom=document.getElementById('general-room')
roomSwitch("general")
