/* 
	This is the main Go file for Empires At War, server edition.

	The server provides a control interface to start a new game.

	As players are added to the game, channels are created to communicate with each player.

	The server must maintain a Websocket on the other end of the player channel.
		Instead of a console-channel-io module, try a websocket-channel-io module
		TODO how would a Websocket server work to be able to connect the channel?
*/


package main


// TODO should there be only one package from which main governs?
// TODO should game be a submodule of controller?
import (
	"log"
	"game"
	"util"
	"fmt"
	"strconv"
)


/*
 * The expected sequence for the game server to get a new game going is:
	
	1. Client calls server to create a new game
		a. server calls game.NewGame; return the game ID to the client
		b. server calls Game.AddPlayer to add the first player
		c. return the game ID to the client
	2. Client calls server with game ID to join game
		a. server calls game.GetGame(gameID) to retrieve game instance
		b. if not found, return error to the client
		c. server calls Game.AddPlayer to add the new player
	3. Client calls server with game ID to start game
		a. server calls game.GetGame(gameID) to retreive game instance
		b. if not found, return error to the client
		c. server calls Game.StartGame

 */



// Run a very simple one-turn game.
func main() {

	util.LoadConfig()	

	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	log.Println("Start eaw")
	
	nbrPlayers := ReadInt("Number of players", 2)
	
	game.InitializeGame(nbrPlayers)
	game.AssignTerritories()
	
	// show initial holdings
	game.PrintTerritories()
	
	// TODO this needs to poll all players and confirm that all want to continue
	for game.ConfirmAllPlayers("Start round?", "n") {
		game.ExecuteRound()
	}

	// document turns
	game.PrintTurns()

	// show final holdings
	game.PrintTerritories()

}


// Use these functions to start the game and round
// TODO


func ReadInt(msg string, def int) int {

	var ans string
	fmt.Printf("%s (%d): ", msg, def)
	_,err := fmt.Scanln(&ans)
	if (err != nil) {
		if (err.Error() == "unexpected newline") {
			return def
		} else {
			panic(err)
		}
	}
	var i int
	i,err = strconv.Atoi(ans)
	if (err != nil) {
		panic(err)
	}
	return i
}
