# Design Doc
This isn't  a document of what has been implemented it is a guide for what I want to be implemented

## Functionality
1. User can join any room(room must have [a-z\-]{3,10} name) using ```/join roomName``` command
1. User can leave any room(room must have [a-z\-]{3,10} name) using ```/leave roomName``` command
1. User can send and recieve message to any room he/she joined
1. User can change handle using the /identify command(later will make usernames unique)


## Data Structures used
1. **Client (struct):**
	```golang
	Client{ //global type
		identifier string
		writeCh *chan Msg // send recieve message from broadcaster
		terminate *chan struct{} // terminate signal
		conn *websocket.Conn
	}
	```

1. **Msg (struct):**
	```golang
	Msg{ //global type
		OpCode  *uint8 `json:"opcode"`
		Content string `json:"content,omitempty"`
		Room    string `json:"room,omitempty"`
		client  *Client //Since the variable is lowercase it isn't marshalled by the json library
		From    *string `json:"from,omitempty"`
	}
	```

1. **Room (map)**
	```golang
	var Room map[*Client]bool //global type
	```

1. **RoomList (map)**
	```golang
	var RoomList map[string]Room //Local to broadcast function
	```

1. **entering (channel)**
	```golang
	var entering := make(chan Msg,chanBuff) //global channel
	```

1. **leaving (channel)**
	```golang
	var leaving := make(chan Msg,chanBuff) //global channel
	```

1. **messaging (channel)**
	```golang
	var messaging:= make(chan Msg,chanBuff) //global channel
	```

## Global Variables(never modified)
We use variables here instead of constants as we want use the address of the following
1. Operation Codes(type is uint8)
	- communicate = 0
	- join = 1
	- leave = 2
	- identify = 3
	- ping = 4
	- leaveAll = 5
1. Default Room Constant
	- defaultRoom = "general"

## Global constants
1. Size constants(type is uint8)
	- clientMsgBuff
	- chanBuff

1. PingTimeout(type is int)

## Functions

### - Incoming connection handler(wsHandler):
**Parameters:** w http.ResponseWriter, r \*http.Request

**Creates:** creates client struct 

**Performs:**
- upgrades connection to websocket
- creates client struct
- sets user identity to random 5 digit number
- starts clientWriter in a goroutine with pointer to client struct as parameter
- enter infinite loop to read user messages, unmarshall the json(user response) into a Msg struct and send it to the appropriate channel based on its opcode
- Its response to opcodes is as follows:
	1. if opcode=0, it adds the user identity to the messages "From" field and sends it over the messaging channel
	1. if opcode=1, it sends the message over the entering channel, then changes the opcode to 0(communicate) and sends it over the messaging channel
	1. if opcode=2 it sends the message over the leaving channel
	1. remaining opcode responses yet to be designed
- exits when socket is closed either by client or due to ReadTimeout(set everytime a message is read succesfully)
- prior to exiting it performs the following cleanup:
	1. closes client.writeCh
	1. send client over leaving channel
	1. send an empty struct over terminate channel(to signal clientWriter to close) and closes the terminate channel
	1. returns to end the goroutine

### - Writing to client(clientWriter)
**Parameters:** cli \*Client

Runs as a goroutine and each Client object has one associated with it

**Creates:**

 a ping Ticker that send a value over a channel every 10 seconds

**Performs:**
- Infinite select statement on the Client.terminate, a ping ticker and Client.writeCh channels
	1. If Client.terminate sends a value it exits
	2. If ping ticker sends a value it sends an empty message to the client with opcode 4
	3. if Client.writeCh sends a value it sends it to the client

### - Broadcast (broadcaster)
**Parameters:** None

Only a single instance of broadcaster is created (as a goroutine) and the roomlist is local to it (to prevent race conditions with RoomList)

**Creates:** RoomList

**Performs:**
- Infinite select over the messaging, entering and leaving channels
	1. If a value is sent on messaging it checks the room and send the message to all clients in that room
	1. If a value is sent on entering it checks with room is specified and adds the client to that room
	1. If a value is sent on leaving it deleted that clients entry in that room and then checks if the room is empty and deletes it if it is
