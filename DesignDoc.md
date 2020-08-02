
# Design Doc

## Functionality
1. User can join any room(roomName must satisfy [a-z0-9\\\-]{2,10} name) using ```/join roomName``` command
1. User can leave any room using ```/leave roomName``` command
1. User can send and receive messages to any room he/she joined
1. User can clear messages in the current room using the ```/clear``` command
1. User can change how he/she is identified using the ```/identify username``` command(username must satisfy '[a-zA-z0-9\\\-]{2,10}')



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
We use variables here instead of constants as we want to use the address of the following
1. Operation Codes(type is uint8)
	- communicate uint8 = 0
	- join        uint8 = 1
	- leave       uint8 = 2
	- identify    uint8 = 3
	- ping        uint8 = 4
	- leaveAll    uint8 = 5
	- notify		uint8 = 6
	- notifyAll	uint8 = 7
1. Default Room Constant
	- defaultRoom = "general"

## Global constants
1. Size constants
	- clientMsgBuff
	- chanBuff

1. PingTimeout

## Functions

### -clientCreator
Create Client object and announce client entering. Runs clientHandler in a goroutine and exits

**Parameters:** w http.ResponseWriter, r \*http.Request

**Creates:** 
- client struct 
- clientHandler go routine 

**Performs:**
- Upgrades connection to websocket
- Creates client struct
- Sets user identity to random 5 digit number
- Announces creation of client through messaging channel
- Adds client to defaultRoom using entering channel
- Creates clientHandler go routine with client as argument and exits

### -clientHandler
Listens for incoming messages from clients and handles them

**Parameters:** client \*Client

**Creates:**
- clientWriter go routine

**Performs:**
- Enter infinite loop to read user messages, unmarshal the JSON(user response) into a Msg struct and send it to ```resolveRequest``` function
- The read deadline for reading user messages is set by ```resolveRequest``` function
- If read error occurs, it runs ```closeClient``` function with and argument of client and exits. 
- If error occurs during umarshaling the incoming message is discarded and we renter the loop

### -resolveRequest

**Parameters:** client \*Client, msg Msg

**Performs:**
- Sets ReadDeadline to ```pingTimeout``` seconds from now
- Sets msg.client to client
- Enters a switch case to handle different opCodes
- Its response to opcodes is as follows:
	1. opcode=0, it adds the user identity to the msg "From" field and sends it over the messaging channel
	1. opcode=1, it notifies users of the room being joined, then adds the user to the room
	1. opcode=2, it removes the user from the group, then notifies other users of the departure 
	1. opcode=3, If the specified username is valid, change user identity and notify users in rooms that the current user has joined


### -clientWriter
Receives messages from cli.writeCh and sends to client. Pings client periodically to prevent ReadDeadline from closing active clients. Runs as a goroutine and each Client object has one associated with it

**Parameters:** cli \*Client

**Creates:**
 a ping Ticker that send a value over a channel every 10 seconds

**Performs:**
- Infinite select statement on the Client.terminate, a ping ticker and Client.writeCh channels
	1. If Client.terminate sends a value it closes ping ticker and exits
	1. If ping ticker sends a value it sends an empty message(ping) to the client with opcode 4
	1. if Client.writeCh sends a value it sends it to the client

### -closeClient
Annouces client departure in all rooms client joined and closes any open channels or goroutines

**Parameters:** client \*Client

**Performs:**
- send client over leaving channel
- send an empty struct over terminate channel(to signal clientWriter to close) and closes the terminate channel
- closes client.writeCh

### -randIdentity

**Parameters:*** None

**Performs:**
- returns a randomly generated 5 digit number

### -usernameValidate
Check if username is valid

***Parameters:*** username string

**Performs:**
- if username has white space it returns false
- if username has non-alphanumeric characters other than _ it returns false
- if username is is longer than 10 characters or shorter than 2 characters it returns false
- else it returns true


### -broadcaster

**Parameters:** None

Only a single instance of broadcaster is created (as a goroutine) and the RoomList is local to it (to prevent race condition)

**Creates:** RoomList

**Performs:**
- Infinite select over the messaging, entering and leaving channels
- If a value(msg) is sent on messaging:
	- if msg.OpCode=7 send a notify opcode to all members in rooms that msg.client has joined
	- else if msg.Room exists, send the msg to all clients in msg.Room  
- If a value(msg) is sent on entering:
	- First, if the specified room doesn't exist in the roomList it creates an empty room
	- Then it adds the client to the room 
- If a value(msg) is sent on leaving:
	1. If opcode=5(leaveAll) it checks every room for msg.client, deletes each match it finds and notifies members of rooms that contained msg.client
	1. If opcode=2(leave) it deletes the client entry in the specified room. No checking is done to see if client is in the room. 
	1. If after the operation is complete the room is empty, the room is deleted from RoomList

### Various message templating functions
Can be found in the client.go module. Used to reduce clutter caused by creating and populating Msg structs in the middle of functions