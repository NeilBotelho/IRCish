identified=false//check if we've identified with server only then allow sendmessage
roomRegex=RegExp('[a-z\-]{3,10}')
//Front end implementation
messageBox=document.getElementById('general-messages')
inputField=document.getElementById('input-field')
help="AVAILABLE COMMANDS:\n"+
"/join channelName : to join a channel(name must contain only lowercase letters and underscores and must begin with a letter)"

function updateMessages(newline,classes=[]){
	var newcontent = document.createElement('p');
	newcontent.innerText = newline;
	while (newcontent.firstChild) {
		   messageBox.appendChild(newcontent.firstChild);
	}
	for(var c in classes){
		newcontent.classList.add(c)
	}
}

function sendMessage(e){
	// Update messageBox
	let msg=e.target.value.trim()
	if(msg.substring(0,5)=="/join"){
		tokens=msg.split(" ")
		if (tokens.length!=2){
			// updateMessages(e.target.value)
			// return
		}
		else{
			if(roomRegex.test(tokens[1])){
				// ws.send(JSON.stringify({"opCode":1,"room":tokens[1]}))
				console.log(tokens[1],tokens.length)
			}
			console.log(msg.split(" ")[1])
			// send join message
		}
		console.log(tokens)	
	}
	else if(msg=="/help"){
		messageBox.innerText+=help
	}
	else{
		if(ws.readyState != WebSocket.OPEN){
			updateMessages("Not connected to server\n")
		}
		else{
    		var newcontent = document.createElement('p');
    		newcontent.innerText = e.target.value+"\n";
    		while (newcontent.firstChild) {
     		   messageBox.appendChild(newcontent.firstChild);
    		}
    		newcontent.id="hello"
			// messageBox.innerText=e.target.value

			//send message over socket
			// ws.send(JSON.stringify({"opCode":2,"content":e.target.value}))
			// console.log()
		}
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
		messageBox.innerText+="\n"+event.data
}
function displayError(){
	alert("Error connecting to server")
}