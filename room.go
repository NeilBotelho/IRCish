package main


const (
	// Channel Buffer Size constants (never changed)
	clientMsgBuff uint8 = 1
	chanBuff      uint8 = 1
	// Ping timeout
	pingTimeout = 15 //seconds
)

var (
	// operations
	communicate uint8 = 0
	join        uint8 = 1
	leave       uint8 = 2
	identify    uint8 = 3
	ping        uint8 = 4
	leaveAll    uint8 = 5
	notify		uint8 = 6
	defaultRoom string = "general"

	//Communication Channel
	entering  = make(chan Msg,chanBuff)
	leaving   = make(chan Msg,chanBuff)
	messaging = make(chan Msg,chanBuff)
)


type Msg struct {
	OpCode  *uint8 `json:"opcode"`
	Content string `json:"content,omitempty"`
	Room    string `json:"room,omitempty"`
	client  *Client
	From    *string `json:"from,omitempty"`
}

type Room map[*Client]bool

func broadcaster() {
	RoomList := make(map[string]Room)
	for {
		select {
		case msg := <-messaging:
			if RoomList[msg.Room] != nil {
				for cli, _ := range RoomList[msg.Room] {
					*cli.writeCh <- msg
				}
			}
		case msg := <-entering:
			if RoomList[msg.Room] == nil {
				RoomList[msg.Room] = Room{}
			}
			RoomList[msg.Room][msg.client] = true

		case msg := <-leaving:
			if *msg.OpCode == leaveAll {
				for roomName, _ := range RoomList {
					delete(RoomList[roomName], msg.client)
					if len(RoomList[roomName]) == 0 {
						delete(RoomList, roomName)
					}
				}
			} else {
				delete(RoomList[msg.Room], msg.client)
				if len(RoomList[msg.Room]) == 0 {
					delete(RoomList, msg.Room)
				}
			}
		}
	}

}