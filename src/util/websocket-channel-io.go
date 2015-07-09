package util

import (
	"fmt"
    . "github.com/gorilla/websocket"
)


// write a string prompt to stdout and read string input from stdin
func WebsocketChannelIO(ch MessageChannel, conn *Conn) {
		
	Mainlog.Println("Started ConsoleChannelIO on channel",ch)
	
	for {
		
		var chanMsg channelMessage
		
		// Mainlog.Println("reading channel",ch)
		chanMsg = <-ch
		
		// TODO define message structure; 
		// start with default JSON rendering of ChannelMessage
		// Use Conn.WriteJSON(chanMsg)
		
		// TODO if the channel message type is MSG_TYPE_REQ, issue a 
		// Conn.ReadJSON()
		// HTML client must send JSON that aligns to ChannelMessage
		

		msg := chanMsg.promptText
		def := chanMsg.defaultText
		
		var ans string
		var err error
		if (chanMsg.messageTypeCode == MSG_TYPE_MSG) || (def == "") {
			_, err = fmt.Printf("%s",msg)
		} else {
			_, err = fmt.Printf("%s (%s): ",msg, def)
		}	
		if (err == nil) {
			if (chanMsg.messageTypeCode == MSG_TYPE_REQ) {
				_,err = fmt.Scanln(&ans)
				if (err != nil) {
					if (err.Error() == "unexpected newline") {
						chanMsg.responseText = def			
					} else {
						chanMsg.err = err
					}
				} else {
					chanMsg.responseText = ans
				}
			}
		} else {
			chanMsg.err = err
		}
		

		// Mainlog.Println("response:",ch,chanMsg)
		ch <- chanMsg
	}
}

