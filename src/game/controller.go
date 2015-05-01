//
package game

import (
	// "log"
	"fmt"
	"container/list"
	"util"
	"strconv"
	"strings"
)

// the maximum number of players the game will allow
var maxNumberOfPlayers = 5

type Game struct {
	Id string
	Players []Player // this will be sized to the maximum number of players supported.
	NumPlayers int // this could be the limit for determining the wrapping of nextPlayerIndex
	Turns *list.List
	CurrentPlayer Player
	Territories []Territory
	TerrMap map[string] *Territory
	PlaySequenceSource int
	nextPlayerIndex int
	sharedConsole bool
}

type play struct {
	sequence int
	attackFrom Territory
	attackTo Territory
	attackerRoll int
	defenderRoll int
	attackerWin bool
}

func (p play) String() string {
	return fmt.Sprint("{",p.sequence,p.attackFrom,p.attackTo,p.attackerRoll,p.defenderRoll,p.attackerWin,"}")
}

type turn struct {
	attacker Player
	plays *list.List
	nbrWins int
	nbrLosses int
	streak int
}

var games map[string] *Game = make(map[string] *Game)

// var players []Player
// var turns *list.List
// var currentPlayer Player
// var playSequenceSource int
// var sharedConsole bool

var gameIdSeq int = 0

func (g Game) GetCurrentPlayer() Player {
	return g.CurrentPlayer
}


func NewGame() *Game {

	// generate a game ID
	gameIdSeq++
	
	// set game ID
	var g Game
	g = Game{}
	g.Id = strconv.Itoa(gameIdSeq)
	
	g.Players = make([]Player, maxNumberOfPlayers)
	g.NumPlayers = 0
	
	// add to game map
	games[g.Id] = &g
	return &g
}

func (g Game) AddPlayer(name string) {
	
	// TODO if g.NumPlayers == maxNumberOfPlayers, we throw an error
	if (g.NumPlayers == maxNumberOfPlayers) {
		util.Mainlog.Println("game.AddPlayer(): game is full of players")
		return
	}
	
	i := g.NumPlayers
	g.NumPlayers++
	
	g.Players[i].UserChan = util.CreateUserChannel()
	// start the goroutine to communicate with this player
	go util.ConsoleChannelIO(g.Players[i].UserChan)
	
	g.Players[i].Name = name
	return
}



func (g Game) InitializeGame() {
	util.Mainlog.Println("game.InitializeGame()")

	// indicate there is only one console being used
	g.sharedConsole = true
	
	// initialize players
	// players = make([]Player, numPlayers)
	
	// for i := 0; i < g.NumPlayers; i++ {
	// 	players[i].UserChan = util.CreateUserChannel()
	// 	// start the goroutine to communicate with this player
	// 	go util.ConsoleChannelIO(players[i].UserChan)
	// }
	//
	// for i := 0; i < numPlayers; i++ {
	// 	// TODO Add prompting for each player's name.
	// 	players[i].Name = "Player " + strconv.Itoa(i+1)
	// }

	util.Mainlog.Println("players: ", g.Players)
	
	
	// Initialize board - getting the board the first time initializes it
	terr := g.LoadTerritories("./etc/terr-in2.json")
	util.Mainlog.Println("Board has", len(terr), "territories")
	
	g.AssignTerritories()
	
	g.PrintTerritories()
	
	// Initialize turns
	g.Turns = list.New()

	// util.Mainlog.Println("turns: ", turns)
	
	util.Init()
}

// Assign territories sequentially 
// TODO this needs to become SelectTerritory such that the players get to pick
func (g Game) AssignTerritories() {
	util.Mainlog.Println("game.AssignTerritories()")
	terr := g.Territories
	
	var j = 0
	for i := range terr {
		// body
		terr[i].Owner = g.Players[j]
		j++
		if j == g.NumPlayers {
			j = 0
		}
		
	}
}

func (g Game) RunGame() {
	
	util.Mainlog.Println("game.RunGame()")
	util.Mainlog.Println("Running game", g)
	
	for g.ConfirmAllPlayers("Start round?", "n") {
		g.ExecuteRound()
	}

	// document turns
	g.PrintTurns()

	// show final holdings
	g.PrintTerritories()
	
	
}

func (g Game) ExecuteRound() {
	
	for pn := 0; pn < len(g.Players); pn++ {
		// start next player turn
		g.StartTurn()
		for g.beginAttackSequence() {
			if (g.ExecutePlay() < 0) {
				break
			}
		}
		g.EndTurn()
	}
}

func (g Game) beginAttackSequence() bool {
	g.printTerritories()
	return g.GetCurrentPlayer().Confirm("Do you want to attack?", "y")
}

// Initialize a new turn for the next player
func (g Game) StartTurn() {
	util.Mainlog.Println("game.StartTurn()")
	g.nextPlayer()
	util.Mainlog.Println("Starting turn for ", g.GetCurrentPlayer())
	g.PutMessageAllPlayers("Starting turn for " + g.GetCurrentPlayer().Name + "\n")
	
	t := turn{g.GetCurrentPlayer(), list.New(), 0, 0, 0}
	g.Turns.PushBack(&t)
}

func (g Game) EndTurn() {

	util.Mainlog.Println("game.EndTurn()")
	g.PutMessageAllPlayers("Ending turn for " + g.CurrentPlayer.Name + "\n")
}

// TODO make players into a ring
// FIXME it should return the next player instead of relying on module var
// var nextPlayerIndex = 0
func (g Game) nextPlayer() {
	util.Mainlog.Println("game.nextPlayer()")
	g.CurrentPlayer = g.Players[g.nextPlayerIndex]
	g.nextPlayerIndex++
	if (g.nextPlayerIndex == len(g.Players)) {
		g.nextPlayerIndex = 0
	}
}

// execute all aspects of one play cycle of the current turn
// return -1 to force end of turn, otherwise return 0
func (g Game) ExecutePlay() int {
	util.Mainlog.Println("game.ExecutePlay()")
	
	// get the current turn from end of turn list
	currTurn := g.Turns.Back().Value.(*turn)
	util.Mainlog.Println("currTurn: ", currTurn)
	
	// create a new play
	var p = play{}
	
	// assign a sequence number
	g.PlaySequenceSource++  // this will need sync'ing in multi-threaded world
	p.sequence = g.PlaySequenceSource
	

	var pDefendTerr *Territory
	var err error
	for {
		// get the attacking territory
		p.attackFrom = g.SelectAttackingTerritory()
		// TODO check for zero-value territory.
		
		// get the defending territory
		pDefendTerr, err = g.SelectDefendingTerritory(p.attackFrom)
		if (err == nil) {
			p.attackTo = *pDefendTerr
			break
		} else {
			g.GetCurrentPlayer().PutMessage(err.Error() + "\n")
		}
	}
	
	// append play to turn's play list - play is "official"
	currTurn.plays.PushBack(&p)

	// roll the die/dice for the attacker
	p.attackerRoll = util.Roll()
	
	// roll the die/dice for the defender
	p.defenderRoll = util.Roll()
	
	g.PutMessageAllPlayers(fmt.Sprintf("Attacker rolled %d, defender rolled %d \n", p.attackerRoll, p.defenderRoll))
	
	
	// determine outcome (win/loss (ties go to defender))
	if (p.attackerRoll > p.defenderRoll) {
		// Current player wins
		g.PutMessageAllPlayers("Attacker wins!\n")
		p.attackerWin = true
		pDefendTerr.Owner = currTurn.attacker

		currTurn.nbrWins++
		if (currTurn.streak < 0) {
			currTurn.streak = 1
		} else {
			currTurn.streak++
		}

	} else {
		// Current player loses
		g.PutMessageAllPlayers("Defender wins!\n")
		currTurn.nbrLosses++
		if (currTurn.streak > 0) {
			currTurn.streak = -1
		} else {
			currTurn.streak--
		}
	}
	
	util.Mainlog.Println("game.executePlay() p:", p)
	
	if (currTurn.nbrLosses == 3) {
		g.PutMessageAllPlayers("That was the third loss in the turn.\n")
		return -1
	} else if (currTurn.streak == -2) {
		g.PutMessageAllPlayers("Two losses in a row. Too bad.\n")
		return -1
	} else {
		return 0
	}
}

// Let the current player select the attacking territory
// from the set of territories owned by the current player
func (g Game) SelectAttackingTerritory() Territory {

	util.Mainlog.Println("game.SelectAttackingTerritory()")
	terr := g.Territories
	var terrIndex int
	var idxMap[] int
	
	// Continue until a valid response is received from the user
	// TODO support a backout option
	for {
		
		idxMap = make([] int, len(terr))
	
		var promptArray []string
		promptArray = make([]string,25)
	
		j := 0
		promptArray[j] = fmt.Sprintln("List of territories owned by", g.GetCurrentPlayer().Name)
		for i := 0; i < len(terr); i++ {
			if (terr[i].Owner == g.GetCurrentPlayer()) {
				j++
				idxMap[j] = i
				promptArray[j+1] = fmt.Sprintf("%d) %s\n", j, terr[i].Name)
			}
		}

		joined := strings.Join(promptArray, "")
		g.GetCurrentPlayer().PutMessage(joined)
		terrIndex = g.GetCurrentPlayer().ReadInt("Select attacking territory: ", 1)
	
		// TODO don't panic, write error message and prompt again
		if (terrIndex <= j) {
			break
		} else {
			g.GetCurrentPlayer().PutMessage(fmt.Sprintf("Invalid entry (%d). Please retry.\n",terrIndex))
		}

	}
	if (terrIndex > 0) {
		// Map the user's input to the actual index in terr
		util.Mainlog.Println("user selected ", terr[idxMap[terrIndex]].Name)
		return terr[idxMap[terrIndex]]
	} else {
		// should only occur if the user enters zero
		return Territory{}
	}
}


// Let the current player select the defending territory
// from the set of territories adjacent to the attacking territory
// NOT owned by the current player.
func (g Game) SelectDefendingTerritory(attackingTerr Territory) (*Territory, error) {

	util.Mainlog.Println("game.SelectDefendingTerritory()")
	terr := g.Territories
	idxMap := make([] int, len(terr))
	
	var promptArray []string
	promptArray = make([]string,25)
	
	j := 0
	promptArray[j] = fmt.Sprintln("List of territories attackable from",attackingTerr.Name)
	for i := 0; i < len(attackingTerr.AttackVectorRefs); i++ {
		if (attackingTerr.AttackVectorRefs[i].Owner != g.GetCurrentPlayer()) {
			j++
			idxMap[j] = i
			promptArray[j+1] = fmt.Sprintf("%d) %s (%s)\n", j, attackingTerr.AttackVectorRefs[i].Name, attackingTerr.AttackVectorRefs[i].Owner.Name)
		}
	}
	// TODO don't panic, require reselection of attacking territory
	if (j == 0) {
		// panic("No territories available to attack.")
		return nil, fmt.Errorf("No territories available to attack from %s.", attackingTerr.Name)
	}

	joined := strings.Join(promptArray, "")
	// fmt.Println(joined)
	g.GetCurrentPlayer().PutMessage(joined)
	terrIndex := g.GetCurrentPlayer().ReadInt("Select territory to attack: ", 1)

	// TODO don't panic, write error message and prompt again
	if (terrIndex > j) {
		panic("User entered invalid input")
	}
	
	util.Mainlog.Println("user selected ", attackingTerr.AttackVectorRefs[idxMap[terrIndex]].Name)
	
	return attackingTerr.AttackVectorRefs[idxMap[terrIndex]], nil
}


func (g Game) PrintTurns() {
	util.Mainlog.Println("game.PrintTurns()")
		
	for et := g.Turns.Front(); et != nil; et = et.Next() {
		t := et.Value.(*turn)
		util.Mainlog.Println("t:",t)
		
		for ep := t.plays.Front(); ep != nil; ep = ep.Next() {
			p := ep.Value.(*play)
			util.Mainlog.Println("p:",p)
		}
	}
	
}


//====================
// short-term archive
//====================
