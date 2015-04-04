package util

import (
	// "fmt"
	// "testing"
	"time"
	"log"
)


func TestConsoleChannelIO() {
	
	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	
	msg := "Here is your prompt: "
	def := "n"
	
	log.Println("TestConsoleChannelIO: start")
	var ch MessageChannel
	
	ch = CreateUserChannel()
	go ConsoleChannelIO(ch)

	log.Println("TestConsoleChannelIO: entering loop")
	
	for {
		var cm channelMessage

		cm.requestTime = time.Now()
		cm.messageTypeCode = MSG_TYPE_REQ
		cm.promptText = msg
		cm.defaultText = def
		log.Println("cm=",cm)
		// convert after channel interaction
		ch <- cm
		cm = <-ch
		cm.responseTime = time.Now()
		log.Println("cm=",cm)
		// t.FailNow()
	}
}