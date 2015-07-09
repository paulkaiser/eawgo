package util

import (
	// "fmt"
	// "log"
	// "game"
)


// write a string prompt to stdout and read string input from stdin
func PlayerChannelIO(p Player) {
	
	Mainlog.Println("Started ChannelIO for player",p)
	if (p.WSConn == nil) {
		ConsoleChannelIO(p.UserChan)
	} else {
		WebsocketChannelIO(p.UserChan, p.WSConn)
	}
}

