package main

// this runs all tests in the util package

import (
	"util"
	// "fmt"
	"log"
)


func main(){
	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	
	log.Println("start testing util package")
	util.TestConsoleChannelIO()
	
}