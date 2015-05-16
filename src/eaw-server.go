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
	To start a new game and become the first player, the client should make the following REST call

	POST http://localhost:8080/eaw/game

	This returns a game ID and a player ID. The client is considered the first player.

	{
		"gameId": "0123456789abcdef",
		"playerId": "0123456789abcdef"
	}

*/
/*
	To join a game, the client should make the following REST call

	POST http://localhost:8080/eaw/game/{game-id}/player

	This returns a player ID that is specific to the game.

	{
		"gameId": "0123456789abcdef",
		"playerId": "0123456789abcdef"
	}

*/
/*
	Once a client has a game ID and player ID, it should open a Websocket and register.

	Registration happens once and causes the server to connect the client Websocket 
	to the player in the game via a channel being monitored in a goroutine.
*/

func main() {
  http.HandleFunc("/echo", echoHandler)
  http.Handle("/", http.FileServer(http.Dir(".")))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("Error: " + err.Error())
  }
}




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
		c. server calls Game.RunGame

 */

/*
	Once a connection is upgraded, it is persistant.

	We can hand the connection over to a function to determine whether
	this connection is from a user session that is already engaged in 
	a game (some kind of session ID, cookie maybe, hmmm can a websocket
	return cookie info?). This kind of disconnect recovery is for later.

	We need a initial handshake protocol. 
		1. Client sends a "New Game" request
		2. Server creates a user session
		3. Server creates a new game and adds user as player
		4. Server starts goroutine passing player channel and websocket conn.
		5. Server returns user session ID and game ID

	The channelMessage type needs to be enhanced.

	Options:
	1.	the channelMessage type has a timestamp, type and payload.
		the payload is JSON

	2.	everything goes JSON

	The websocket channel IO controller will need to convert
	the channelMessage to JSON and translate the response,
	if required, to from JSON to a channelMessage.


	Need more robust message protocol and handling.
	Start with the two handshakes we currently have:
		1. Message - server sends information to client for display
		2. Request - server sends a specific request to the client and 
			expects a specific response.
*/



func connectionHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        panic("Error: " + err.Error())
    }
 
 	fmt.Println("conn: ",conn)

	// TODO check if there is already a user session on the connection
	
	// No session - attach the connection to a game
	// TODO assign a session to the connection
	
}


// Run a very simple one-turn game.
func main0() {

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
	// g.RunGame()
	
	// We can't RunGame() in a goroutine until we know that the main function will
	// not exit.
	go g.RunGame()
	// return acknowledgement
	
	
	// is there a sleep function in go
	select {}
	
}


