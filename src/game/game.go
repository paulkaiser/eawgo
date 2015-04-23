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
// var territories []Territory

// indexes territories by Name
// var terrMap map[string] *Territory

func (g Game) GetTerritories() []Territory {
	return g.Territories
}
func (g Game) PrintTerritories() {
	logTerritories(g.GetTerritories())
	g.printTerritories()
}

func (g Game) printTerritories() {
	util.Mainlog.Println("game.printTerritories()")
	
	for i := 0; i < len(g.GetTerritories()); i++ {
		g.PutMessageAllPlayers("Territory: " + g.GetTerritories()[i].ShortDescription() + "\n")
	}
}
func logTerritories(territories []Territory) {
	util.Mainlog.Println("game.logTerritories()")
	
	for i := 0; i < len(territories); i++ {
		util.Mainlog.Println("Territory:", territories[i])
	}
}

func (g Game) LoadTerritories(filename string) []Territory {
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
	g.mapTerritories()
	g.generateAttackVectors()
	return territories
}

func (g Game) mapTerritories() {
	g.TerrMap = make(map[string] *Territory)
	for i := 0; i < len(g.Territories); i++ {
		terrMap[g.Territories[i].Name] = &g.Territories[i]
	}
	
	util.Mainlog.Println("terrMap:", g.TerrMap)
}

// Take the array of string AttackVectors and
// convert it to the array of *Territory AttackVectorRefs
func (g Game) generateAttackVectors() {

	for i := 0; i < len(g.Territories); i++ {
		g.Territories[i].AttackVectorRefs = make([] *Territory, len(g.Territories[i].AttackVectors))
		for j := 0; j < len(g.Territories[i].AttackVectors); j++ {
			g.Territories[i].AttackVectorRefs[j] = g.TerrMap[g.Territories[i].AttackVectors[j]]
		}
	}
	util.Mainlog.Println("terr:",g.Territories)
}

//
func (g Game) SaveTerritories(filename string) {
	util.Mainlog.Println("game.SaveTerritories(",filename,")")
	
	// save list of territories in JSON format
	terr := g.GetTerritories()
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