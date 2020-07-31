package main

const (
	// Channel Buffer Size constants (never changed)
	clientMsgBuff uint8 = 1
	chanBuff      uint8 = 1
	// Ping timeout
	pingTimeout = 25 //seconds
)


var (
	// Operation codes(opCodes used in Msg)
	// These are never changed
	//  We make them vars instead of const as we want to address them 
	communicate uint8 = 0
	join        uint8 = 1
	leave       uint8 = 2
	identify    uint8 = 3
	ping        uint8 = 4
	leaveAll    uint8 = 5
	notify		uint8 = 6
	notifyAll	uint8 = 7
	defaultRoom string = "general"

	//Broadcast Channels
	entering  = make(chan Msg,chanBuff)
	leaving   = make(chan Msg,chanBuff)
	messaging = make(chan Msg,chanBuff)
)


// Represents a message to or from client
// Not all fields will be used in every operation. 
type Msg struct {
	OpCode  *uint8 `json:"opcode"` 
	Content string `json:"content,omitempty"`
	Room    string `json:"room,omitempty"`
	client  *Client
	From    *string `json:"from,omitempty"`
}


// Create Room type
// using map for fast lookups even though the value it returns is arbitrary
type Room map[*Client]bool


func broadcaster() {
/*Creates and controls access to RoomList
All join, leave and message operations occur through this function*/	
	
	// Map of roomName to Room
	RoomList := make(map[string]Room)
	for {
		select {
		case msg := <-messaging:
			if *msg.OpCode == notifyAll{
				// send a notify opcode to all 
				// members in rooms that msg.client has joined				
				msg.OpCode=&notify
				for roomName, _ := range RoomList {
					if(RoomList[roomName][msg.client]){
						for cli, _ := range RoomList[roomName] {
							msg.Room=roomName
							*cli.writeCh <- msg
						}
					}
				}
			}else if RoomList[msg.Room] != nil{
				// Here *msg.OpCode == communication
				// send the msg to all clients in msg.Room
				for cli, _ := range RoomList[msg.Room] {
					*cli.writeCh <- msg
				}
			}
		case msg := <-entering:
			// Adds client to msg.Room
			// creates the room first if it doesn't exist
			if RoomList[msg.Room] == nil {
				RoomList[msg.Room] = Room{}
			}
			RoomList[msg.Room][msg.client] = true

		case msg := <-leaving:
			if *msg.OpCode == leaveAll {
				// Remove client from all rooms and notify room members
				
				msg.OpCode=&notify
				// iterate over all rooms 
				for roomName, _ := range RoomList {
					// if msg.client is in the room
					if(RoomList[roomName][msg.client]){
						// delete msg.client from the room
						delete(RoomList[roomName], msg.client)
						for cli, _ := range RoomList[roomName] {
							// Notify members of the room of msg.client's departure
							msg.Room=roomName
							*cli.writeCh <- msg
						}
					}
					if len(RoomList[roomName]) == 0 {
						// If room is empty, remove room from RoomList
						delete(RoomList, roomName)
					}
				}
			} else {
				// Here *msg.OpCode == leave 
				// Remove client from 
				delete(RoomList[msg.Room], msg.client)
				if len(RoomList[msg.Room]) == 0 {
					delete(RoomList, msg.Room)
				}
			}
		}
	}

}
