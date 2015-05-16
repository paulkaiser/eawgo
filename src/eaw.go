/* 
	This is the main Go file for Empires At War.
	Initially, it is intended to capture:

		Overall game structure that is not yet expressed in sub-modules
		Utility modules not yet coded
		Changes to existing modules not yet coded and/or expressed in TODO or FIXME tags
*/


package main


// TODO should there be only one package from which main governs?
// TODO should game be a submodule of controller?
import (
	"log"
	"game"
	"util"
	// "fmt"
	// "strconv"
)


// TODO need to play with user keyboard input before moving to a client-server model.
// TODO need to think about ending the turn (3 total losses or 2 consecutive losses or by choice)
// TODO how to end the game?

// TODO use channels to gather input from users and send to controller

// Run a very simple one-turn game.
func main() {


	util.LoadConfig()	

	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	log.Println("Start eaw-server [test]")

	// client calls to create new game
	var g *game.Game
	g = game.NewGame()
	g.AddPlayer("Player 1")
	// return game ID to caller
	
	// client calls to join game using ID
	g.AddPlayer("Player 2")
	// return acknowledgement
	
	// client called to start the game
	g.InitializeGame()
	g.RunGame()
	
}


// Use these functions to start the game and round
// TODO


// func ReadInt(msg string, def int) int {
//
// 	var ans string
// 	fmt.Printf("%s (%d): ", msg, def)
// 	_,err := fmt.Scanln(&ans)
// 	if (err != nil) {
// 		if (err.Error() == "unexpected newline") {
// 			return def
// 		} else {
// 			panic(err)
// 		}
// 	}
// 	var i int
// 	i,err = strconv.Atoi(ans)
// 	if (err != nil) {
// 		panic(err)
// 	}
// 	return i
// }
