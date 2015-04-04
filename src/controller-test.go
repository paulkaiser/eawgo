// test the dice rolling package
package main

import (
	"log"
	"controller"
	"game"
	"fmt"
)

func main() {

	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	log.Println("Start controller-test")
	
	controller.InitializeGame(2)
	controller.AssignTerritories()
	
	// show initial holdings
	game.PrintTerritories()
	
	// player 1 turn
	controller.StartTurn()
	for confirm("Do you want to attack?") {
		controller.ExecutePlay()
		
	}
	controller.EndTurn()

	// player 2 turn
	controller.StartTurn()
	
	for confirm("Do you want to attack?") {
		controller.ExecutePlay()
		
	}
	controller.EndTurn()

	// document turns
	controller.PrintTurns()

	// show final holdings
	game.PrintTerritories()

}

func confirm(msg string) bool {
	
	var ans string
	fmt.Printf("%s (y/n): ",msg)
	_,err := fmt.Scanln(&ans)
	if (err != nil) {
		panic(err)
	}
	if ans == "y" {
		return true
	} else {
		return false
	}
}
