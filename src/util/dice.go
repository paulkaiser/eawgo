// a dice rolling package
package util

import (
	// "log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	defaultSeedValue = 41
	seedVarName = "EAW_DICE_SEED"
)

/*
	Initialize the dice package.
	Seeds the pseudo-random number generator from EAW_DICE_SEED or from default if not provided
*/
func Init() {
	
	var seedValue int64

	val := os.Getenv(seedVarName)
	tmp,err := strconv.Atoi(val)
	if (err != nil) {
		seedValue = generateSeed() //defaultSeedValue
	} else {
		seedValue = int64(tmp)
	}
	rand.Seed(seedValue)
	Mainlog.Println("Initialized rand with seed of",seedValue)
}

/*
	Rolls one 6-sided die and returns the result
*/
func Roll() int {

//	var numRoll int = 1

	var roll int
	roll = rand.Intn(6) + 1
	return roll;
}

// generate a seed value from time components
func generateSeed() int64 {
	
	n := time.Now()
	i := int64(n.Nanosecond())
	return i
}


/*	Consider this as a way to get a better random number seed, 
	a prime number <= a given integer.

// A concurrent prime sieve

package main

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	for i := 0; i < 100; i++ {
		prime := <-ch
		print(prime, "\n")
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}

*/