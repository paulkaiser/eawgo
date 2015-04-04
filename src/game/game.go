//
package game

import (
	// "log"
	"sync"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"util"
)


type Territory struct {
	// The name of the territory. It must be unique
	Name string
	// The player that currently controls the territory
	Owner Player
	// The territories that can be attacked from this territory
	AttackVectors[] string
	AttackVectorRefs[] *Territory
	// The JSON may become a CSV so strings can be converted to references as needed
	// without requiring redundant storage.
}

func (t Territory) String() string {
	return fmt.Sprint("{Name=",t.Name," Owner.Name=",t.Owner.Name," AttackVectors=",t.AttackVectors,"}")	
}
func (t Territory) ShortDescription() string {
	return fmt.Sprint("{Name=",t.Name," Owner.Name=",t.Owner.Name,"}")	
}

var once sync.Once

// stores territories
var territories []Territory

// indexes territories by Name
var terrMap map[string] *Territory

func GetTerritories() []Territory {
	return territories
}
func PrintTerritories() {
	logTerritories()
	printTerritories()
}

func printTerritories() {
	util.Mainlog.Println("game.printTerritories()")
	
	for i := 0; i < len(territories); i++ {
		PutMessageAllPlayers("Territory: " + territories[i].ShortDescription() + "\n")
	}
}
func logTerritories() {
	util.Mainlog.Println("game.logTerritories()")
	
	for i := 0; i < len(territories); i++ {
		util.Mainlog.Println("Territory:", territories[i])
	}
}

func LoadTerritories(filename string) []Territory {
	util.Mainlog.Println("game.LoadTerritories(",filename,")")
	
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	
	err = json.Unmarshal(b, &territories)
	if (err != nil) {
		panic(err)
	}
	// util.Mainlog.Println("terr:",territories)
	mapTerritories()
	generateAttackVectors()
	return territories
}

func mapTerritories() {
	terrMap = make(map[string] *Territory)
	for i := 0; i < len(territories); i++ {
		terrMap[territories[i].Name] = &territories[i]
	}
	
	util.Mainlog.Println("terrMap:", terrMap)
}

// Take the array of string AttackVectors and
// convert it to the array of *Territory AttackVectorRefs
func generateAttackVectors() {

	for i := 0; i < len(territories); i++ {
		territories[i].AttackVectorRefs = make([] *Territory, len(territories[i].AttackVectors))
		for j := 0; j < len(territories[i].AttackVectors); j++ {
			territories[i].AttackVectorRefs[j] = terrMap[territories[i].AttackVectors[j]]
		}
	}
	util.Mainlog.Println("terr:",territories)
}

//
func SaveTerritories(filename string) {
	util.Mainlog.Println("game.LoadTerritories(",filename,")")
	
	// save list of territories in JSON format
	terr := territories
	b, err := json.Marshal(terr)
	if err != nil {
		panic(err)
	}
	s := string(b)
	util.Mainlog.Println("s:", s)
	
	err = ioutil.WriteFile(filename, b, 0644)
	if (err != nil) {
		panic(err)
	}
}