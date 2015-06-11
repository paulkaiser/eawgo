package game

import (
	"util"
	"fmt"
	// "log"
	"strings"
	"github.com/gorilla/websocket"
)

type Player struct {
	Name string
	
	// The Game will use this channel to communicate to the player
	UserChan util.MessageChannel  // try to hide channel type

	// for use in HTML client. If this is nil when the channel IO handler is started,
	// the handler assumes use of the console.
	WSConn *Conn
}

func (p Player) String() string {
	return fmt.Sprint("{",p.Name,"}")
}


func (p Player) Confirm(msg string, def string) bool {

	var ans bool

	pa := []string{p.Name, msg}
	// ans = util.Confirm(msg, def)
	ans = util.GetChannelBoolean(p.UserChan, strings.Join(pa, ", "), def)
	// log.Println("ans=",ans)
	// handle err?
	
	return ans
}

func (g Game) ConfirmAllPlayers(msg string, def string) bool {
	
	confirmedAll := make([] bool, g.NumPlayers)
	
	// Since this is intended to poll all players, do not short-circuit.
	for i := range confirmedAll {
		confirmedAll[i] = g.Players[i].Confirm(msg, def)
	}

	var confirmed bool = true
	for i := range confirmedAll {
		confirmed = (confirmed && confirmedAll[i])
	}
	
	return confirmed
}



func (p Player) ReadString(msg string, def string) string {

	var ans string
	// var err error
	
	pa := []string{p.Name, msg}
	// ans = util.ReadString(msg, def)
	ans = util.GetChannelString(p.UserChan, strings.Join(pa, ", "), def)
	// handle err?
	
	return ans
}

func (p Player) ReadInt(msg string, def int) int {

	var ans int
	// var err error
	
	pa := []string{p.Name, msg}
	// ans = util.ReadInt(msg, def)
	ans = util.GetChannelInteger(p.UserChan, strings.Join(pa, ", "), def)
	// handle err?

	return ans
}

func (p Player) PutMessage(msg string) {
	
	util.PutChannelMessage(p.UserChan, msg)
	// util.PutMessage(msg)
}

// broadcast a message such that all players will "see" it
func (g Game) PutMessageAllPlayers(msg string)  {
	
	if g.sharedConsole {
		
		// Since only one console is in use, use player 1 console
		g.Players[0].PutMessage(msg)
		
	} else {
	
		// Since this is intended to reach all players, do not short-circuit.
		for i := range g.Players {
			g.Players[i].PutMessage(msg)
		}
	}
}

