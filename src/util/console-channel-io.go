package util

import (
	"fmt"
	// "log"
)


// write a string prompt to stdout and read string input from stdin
func ConsoleChannelIO(ch MessageChannel) {
	
	Mainlog.Println("Started ConsoleChannelIO on channel",ch)
	
	for {
		
		var chanMsg channelMessage
		
		// Mainlog.Println("reading channel",ch)
		chanMsg = <-ch
		
		// TODO put log level wrapper around log package
		// Mainlog.Println("request:",chanMsg)
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

