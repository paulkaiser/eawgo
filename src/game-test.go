// test the dice rolling package
package main

import (
	"log"
	"game"
)

func main() {

	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	log.Println("Start game-test")
	
	game.LoadTerritories("./etc/terr-in2.json")
	terr := game.GetTerritories()
	log.Println("len(terr) =",len(terr))
	
	log.Println("End game-test")
}