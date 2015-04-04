// test the dice rolling package
package main

import (
	"log"
	"controller/dice"
)

func main() {

	dice.Init()
	for i := 0; i < 1000; i++ {
		r := dice.Roll()
		if (r < 1) || (r > 6) {
			log.Fatalf("Bad roll outcome: %d\n",r)
		}
		log.Printf("rolled %d\n", r)
	}

}

