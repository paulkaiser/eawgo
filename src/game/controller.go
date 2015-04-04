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

type game struct {
	id string
	players []Player
	turns *list.List
	currentPlayer Player
	playSequenceSource int
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

var games map[string] *game

var players []Player
var turns *list.List
var currentPlayer Player
var playSequenceSource int
var sharedConsole bool

func GetCurrentPlayer() Player {
	return currentPlayer
}

func InitializeGame(numPlayers int) {
	util.Mainlog.Println("game.InitializeGame()")

	// indicate there is only one console being used
	sharedConsole = true
	
	// initialize players
	players = make([]Player, numPlayers)
	
	for i := 0; i < numPlayers; i++ {
		players[i].UserChan = util.CreateUserChannel()
		// start the goroutine to communicate with this player
		go util.ConsoleChannelIO(players[i].UserChan)
	}
	
	for i := 0; i < numPlayers; i++ {
		// TODO Add prompting for each player's name.
		players[i].Name = "Player " + strconv.Itoa(i+1)
	}

	util.Mainlog.Println("players: ", players)
	
	
	// Initialize board - getting the board the first time initializes it
	terr := LoadTerritories("./etc/terr-in2.json")
	util.Mainlog.Println("Board has", len(terr), "territories")
	
	// Initialize turns
	turns = list.New()
	// util.Mainlog.Println("turns: ", turns)
	
	util.Init()
}

// Assign territories sequentially 
// TODO this needs to become SelectTerritory such that the players get to pick
func AssignTerritories() {
	util.Mainlog.Println("game.AssignTerritories()")
	terr := GetTerritories()
	
	var j = 0
	for i := range terr {
		// body
		terr[i].Owner = players[j]
		j++
		if j == len(players) {
			j = 0
		}
		
	}
}

func ExecuteRound() {
	
	for pn := 0; pn < len(players); pn++ {
		// start next player turn
		StartTurn()
		for beginAttackSequence() {
			if (ExecutePlay() < 0) {
				break
			}
		}
		EndTurn()
	}
}

func beginAttackSequence() bool {
	printTerritories()
	return GetCurrentPlayer().Confirm("Do you want to attack?", "y")
}

// Initialize a new turn for the next player
func StartTurn() {
	util.Mainlog.Println("game.StartTurn()")
	nextPlayer()
	util.Mainlog.Println("Starting turn for ", currentPlayer)
	PutMessageAllPlayers("Starting turn for " + currentPlayer.Name + "\n")
	
	t := turn{currentPlayer, list.New(), 0, 0, 0}
	turns.PushBack(&t)
}

func EndTurn() {

	util.Mainlog.Println("game.EndTurn()")
	PutMessageAllPlayers("Ending turn for " + currentPlayer.Name + "\n")
}

// TODO make players into a ring
// FIXME it should return the next player instead of relying on module var
var nextPlayerIndex = 0
func nextPlayer() {
	util.Mainlog.Println("game.nextPlayer()")
	currentPlayer = players[nextPlayerIndex]
	nextPlayerIndex++
	if (nextPlayerIndex == len(players)) {
		nextPlayerIndex = 0
	}
}

// execute all aspects of one play cycle of the current turn
// return -1 to force end of turn, otherwise return 0
func ExecutePlay() int {
	util.Mainlog.Println("game.ExecutePlay()")
	
	// get the current turn from end of turn list
	currTurn := turns.Back().Value.(*turn)
	util.Mainlog.Println("currTurn: ", currTurn)
	
	// create a new play
	var p = play{}
	
	// assign a sequence number
	playSequenceSource++  // this will need sync'ing in multi-threaded world
	p.sequence = playSequenceSource
	

	var pDefendTerr *Territory
	var err error
	for {
		// get the attacking territory
		p.attackFrom = SelectAttackingTerritory()
		// TODO check for zero-value territory.
		
		// get the defending territory
		pDefendTerr, err = SelectDefendingTerritory(p.attackFrom)
		if (err == nil) {
			p.attackTo = *pDefendTerr
			break
		} else {
			currentPlayer.PutMessage(err.Error() + "\n")
		}
	}
	
	// append play to turn's play list - play is "official"
	currTurn.plays.PushBack(&p)

	// roll the die/dice for the attacker
	p.attackerRoll = util.Roll()
	
	// roll the die/dice for the defender
	p.defenderRoll = util.Roll()
	
	PutMessageAllPlayers(fmt.Sprintf("Attacker rolled %d, defender rolled %d \n", p.attackerRoll, p.defenderRoll))
	
	
	// determine outcome (win/loss (ties go to defender))
	if (p.attackerRoll > p.defenderRoll) {
		// Current player wins
		PutMessageAllPlayers("Attacker wins!\n")
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
		PutMessageAllPlayers("Defender wins!\n")
		currTurn.nbrLosses++
		if (currTurn.streak > 0) {
			currTurn.streak = -1
		} else {
			currTurn.streak--
		}
	}
	
	util.Mainlog.Println("game.executePlay() p:", p)
	
	if (currTurn.nbrLosses == 3) {
		PutMessageAllPlayers("That was the third loss in the turn.\n")
		return -1
	} else if (currTurn.streak == -2) {
		PutMessageAllPlayers("Two losses in a row. Too bad.\n")
		return -1
	} else {
		return 0
	}
}

// Let the current player select the attacking territory
// from the set of territories owned by the current player
func SelectAttackingTerritory() Territory {

	util.Mainlog.Println("game.SelectAttackingTerritory()")
	terr := GetTerritories()
	var terrIndex int
	var idxMap[] int
	
	// Continue until a valid response is received from the user
	// TODO support a backout option
	for {
		
		idxMap = make([] int, len(terr))
	
		var promptArray []string
		promptArray = make([]string,25)
	
		j := 0
		promptArray[j] = fmt.Sprintln("List of territories owned by",currentPlayer.Name)
		for i := 0; i < len(terr); i++ {
			if (terr[i].Owner == currentPlayer) {
				j++
				idxMap[j] = i
				promptArray[j+1] = fmt.Sprintf("%d) %s\n", j, terr[i].Name)
			}
		}

		joined := strings.Join(promptArray, "")
		currentPlayer.PutMessage(joined)
		terrIndex = currentPlayer.ReadInt("Select attacking territory: ", 1)
	
		// TODO don't panic, write error message and prompt again
		if (terrIndex <= j) {
			break
		} else {
			currentPlayer.PutMessage(fmt.Sprintf("Invalid entry (%d). Please retry.\n",terrIndex))
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
func SelectDefendingTerritory(attackingTerr Territory) (*Territory, error) {

	util.Mainlog.Println("game.SelectDefendingTerritory()")
	terr := GetTerritories()
	idxMap := make([] int, len(terr))
	
	var promptArray []string
	promptArray = make([]string,25)
	
	j := 0
	promptArray[j] = fmt.Sprintln("List of territories attackable from",attackingTerr.Name)
	for i := 0; i < len(attackingTerr.AttackVectorRefs); i++ {
		if (attackingTerr.AttackVectorRefs[i].Owner != currentPlayer) {
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
	currentPlayer.PutMessage(joined)
	terrIndex := currentPlayer.ReadInt("Select territory to attack: ", 1)

	// TODO don't panic, write error message and prompt again
	if (terrIndex > j) {
		panic("User entered invalid input")
	}
	
	util.Mainlog.Println("user selected ", attackingTerr.AttackVectorRefs[idxMap[terrIndex]].Name)
	
	return attackingTerr.AttackVectorRefs[idxMap[terrIndex]], nil
}


func PrintTurns() {
	util.Mainlog.Println("game.PrintTurns()")
		
	for et := turns.Front(); et != nil; et = et.Next() {
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
