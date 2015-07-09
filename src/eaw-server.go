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
    "github.com/gorilla/websocket"
	"net/http"
	"time"
	"fmt"
	// "strconv"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
 

/*
	To start a new game and become the first player, the user should point their browser to 

	GET http://localhost:8080/eaw/game

	This starts a new game and returns the game ID and a player ID. 
	The response is a page that loads the game client.
	The client is considered the first player.

	{
		"gameId": "0123456789abcdef",
		"playerId": "0123456789abcdef"
	}

	The response also sets a session cookie that is specific to the game and player.

*/
/*
	To join a game, the user points their browser to

	GET http://localhost:8080/eaw/player?game={game-id}

	This creates a new player and a player ID that is specific to the game.
	The response is a page that loads the game client.
	The client is considered the new player.

	{
		"gameId": "0123456789abcdef",
		"playerId": "0123456789abcdef"
	}

	The response also sets a session cookie that is specific to the game a player.

*/
/*
	Once a client has a game ID and player ID, it should open a Websocket and register.

	ws://localhost:8080/eaw/ws
	
	Registration happens once and causes the server to connect the client Websocket 
	to the player in the game via a channel being monitored in a goroutine.

	The request should include the session cookie.

*/

func main() {

	util.LoadConfig()	

	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	log.Println("Start eaw-server [test]")


	http.HandleFunc("/eaw/game", newGameHandler)
	http.HandleFunc("/eaw/player", addPlayerHandler)
	http.HandleFunc("/eaw/ws", wsConnectHandler)

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

// TODO remove this hack
var gameId string


func newGameHandler(w http.ResponseWriter, r *http.Request) {
 
 	fmt.Println("newGameHandler()")
 	fmt.Println("r: ",r)

	// TODO check if there is already a user session on the connection
	
	// No session - attach the connection to a game
	// TODO assign a session to the connection
	
	
	// Create new game and add player
	var g *game.Game
	g = game.NewGame()
	g.AddPlayer("JoJo")
	// return game ID to caller
	gameId = g.Id
	

	// write a short-term cookie
	expiration := time.Now().Add(1 * 24 * time.Hour)
	cookie := http.Cookie{Name: "eaw-session", Value: "game="+g.Id+",player=JoJo", Expires: expiration}
	http.SetCookie(w, &cookie)
	
	
	// TEST redirect to /eaw-client.html
	// TODO firgure out how to get game ID back to client
	
	http.Redirect(w, r, "/eaw-client.html", http.StatusFound)

 	fmt.Println("w: ",w)
}


func addPlayerHandler(w http.ResponseWriter, r *http.Request) {
 
 	fmt.Println("addPlayerHandler()")
 	fmt.Println("r: ",r)

	err := r.ParseForm()
    if err != nil {
        panic("Error: " + err.Error())
    }

	fmt.Println("r.Form: ", r.Form)
	
	// TODO check if there is already a user session on the connection
	cookie, _ := r.Cookie("eaw-session")
	fmt.Println("cookie: ", cookie)
	
	// No session - attach the connection to a game
	// TODO get player name into the request
	// TODO extract the player name from the request
	
	// TODO assign a session to the connection
	name := "KoKo"
	expiration := time.Now().Add(1 * 24 * time.Hour)
	eawCookie := http.Cookie{Name: "eaw-session", Value: "game="+gameId+",player="+name, 
		Expires: expiration}
	http.SetCookie(w, &eawCookie)
	
	// TEST redirect to /eaw-client.html
	http.Redirect(w, r, "/eaw-client.html", http.StatusFound)

 	fmt.Println("w: ",w)
}


func wsConnectHandler(w http.ResponseWriter, r *http.Request) {
	
	var Conn *conn
	var err error
	
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        panic("Error: " + err.Error())
    }
 
 	fmt.Println("wsConnectHandler()")
 	fmt.Println("r: ",r)
	cookie, _ := r.Cookie("eaw-session")
	fmt.Println("cookie: ", cookie)
 	fmt.Println("w: ",w)
 	fmt.Println("conn: ",conn)
	
	// TODO use eaw-session to get handle to game and player
	
	// TODO extract session data not available from conn
	// and pass to keepListening
	go keepListening(conn)

	// TODO check if there is already a user session on the connection
	
	// No session - attach the connection to a game
	// TODO assign a session to the connection
	
}


func keepListening(conn *websocket.Conn) {
	
	// loop forever
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            return
        }
		
		switch (messageType) {
			
		case websocket.BinaryMessage:
			print_binary(p)
		
		case websocket.TextMessage:
			fmt.Printf("Received text: %s\n", string(p))
			
		default:
			fmt.Printf("Received message type %d\n", messageType)
			
		}
		
        // print_binary(p)
 
        err = conn.WriteMessage(messageType, p);
        if  err != nil {
            return
        }
    }
	
}


func print_binary(s []byte) {
  fmt.Printf("Received b:");
  for n := 0;n < len(s);n++ {
    fmt.Printf("%d,",s[n]);
  }
  fmt.Printf("\n");
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


