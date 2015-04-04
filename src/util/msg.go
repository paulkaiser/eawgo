package util

import (
	"time"
)


const (
	MSG_TYPE_MSG = "M"
	MSG_TYPE_REQ = "R"
)
/* 
	This is a private structure used to communicate between the util
	functions that return a specific type and the channel
*/
// TODO this should be private, but is exported to facilitate testing
type channelMessage struct {
	
	//
	requestTime time.Time
	// promptCode string
	messageTypeCode string
	promptText string
	defaultText string

	//
	responseTime time.Time
	responseText string
	
	//
	err error
}

type MessageChannel chan channelMessage

