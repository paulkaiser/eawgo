package util

import (
	"fmt"
	"strconv"
	"time"
	// "log"
)

func CreateUserChannel() MessageChannel {
	
	ch := make(MessageChannel)
	return ch
}

//==== channel-specific functions
//==== read msg from channel
func GetChannelBoolean(ch MessageChannel, msg string, def string) bool {

	var retVal bool

	var cm channelMessage
	cm.requestTime = time.Now()
	cm.messageTypeCode = MSG_TYPE_REQ
	cm.promptText = msg
	cm.defaultText = def
	// util.Mainlog.Println("cm=",cm)
	// convert after channel interaction
	ch <- cm
	cm = <-ch
	cm.responseTime = time.Now()
	Mainlog.Println("cm=",cm)
	// TODO check error
	retVal = convertStringToBool(cm.responseText)
	return retVal
}

func GetChannelInteger(ch MessageChannel, msg string, def int) int {
	
	var retVal int

	var cm channelMessage
	cm.requestTime = time.Now()
	cm.messageTypeCode = MSG_TYPE_REQ
	cm.promptText = msg
	cm.defaultText = strconv.Itoa(def)
	// util.Mainlog.Println("cm=",cm)
	// convert after channel interaction
	ch <- cm
	cm = <-ch
	cm.responseTime = time.Now()
	// util.Mainlog.Println("cm=",cm)
	// TODO check error
	retVal = convertStringToInt(cm.responseText)
	return retVal
}

func GetChannelString(ch MessageChannel, msg string, def string) string {

	var retVal string

	var cm channelMessage
	cm.requestTime = time.Now()
	cm.messageTypeCode = MSG_TYPE_REQ
	cm.promptText = msg
	cm.defaultText = def
	// Mainlog.Println("ch=",ch)
	// Mainlog.Println("cm=",cm)
	// convert after channel interaction
	ch <- cm
	cm = <-ch
	cm.responseTime = time.Now()
	// util.Mainlog.Println("cm=",cm)
	// TODO check error
	retVal = cm.responseText
	return retVal
	
	return retVal
}

func PutChannelMessage(ch MessageChannel, msg string) {

	// fmt.Printf("%s\n", msg)
	
	var cm channelMessage
	cm.requestTime = time.Now()
	cm.messageTypeCode = MSG_TYPE_MSG
	cm.promptText = msg
	// cm.defaultText = strconv.Itoa(def)
	// util.Mainlog.Println("cm=",cm)
	// convert after channel interaction
	ch <- cm
	cm = <-ch
	// TODO check for cm.err
}

// 
// func PutDistinctChannelMessage(ch []MessageChannel, msg string) {
//
// 	// fmt.Printf("%s\n", msg)
//
// 	var cm channelMessage
// 	cm.requestTime = time.Now()
// 	cm.messageTypeCode = MSG_TYPE_MSG
// 	cm.promptText = msg
// 	// cm.defaultText = strconv.Itoa(def)
// 	// util.Mainlog.Println("cm=",cm)
// 	// convert after channel interaction
// 	ch <- cm
// 	cm = <-ch
// 	// TODO check for cm.err
// }


//=== private conversion functions
func convertStringToBool(ans string) bool {
	
	if ans == "y" {
		return true
	} else {
		return false
	}
}

func convertStringToInt(ans string) int {

	i,err := strconv.Atoi(ans)
	if (err != nil) {
		panic(err)
	}
	return i
}

// = archive
// func Confirm(msg string, def string) bool {
//
// 	var ans string
// 	fmt.Printf("%s (%s): ",msg, def)
// 	_,err := fmt.Scanln(&ans)
// 	if (err != nil) {
// 		if (err.Error() == "unexpected newline") {
// 			ans = def
// 		} else {
// 			panic(err)
// 		}
// 	}
// 	if ans == "y" {
// 		return true
// 	} else {
// 		return false
// 	}
// }
//
// func ReadString(msg string, def string) string {
//
// 	var ans string
// 	fmt.Printf("%s (%s): ",msg, def)
// 	_,err := fmt.Scanln(&ans)
// 	if (err != nil) {
// 		if (err.Error() == "unexpected newline") {
// 			ans = def
// 		} else {
// 			panic(err)
// 		}
// 	}
// 	return ans
// }
//
//
func PutMessage(msg string) {

	fmt.Printf("%s\n", msg)

}

