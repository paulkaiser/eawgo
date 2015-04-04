// config functions
package util

import (
	// "encoding/json"
	// "io/ioutil"
	"log"
	// "fmt"
	"os"
	"syscall"
)


var Stdout	*log.Logger
var Stderr	*log.Logger
var Mainlog *log.Logger

// Load the configuration file and expose 
func LoadConfig() {
	
	flags := log.Ldate | log.Lmicroseconds | log.Lshortfile

	// create default loggers
	Stdout = log.New(os.NewFile(uintptr(syscall.Stdout), "/dev/stdout"), "", flags)
	Stderr = log.New(os.NewFile(uintptr(syscall.Stderr), "/dev/stderr"), "", flags)

	filename := "./etc/config.json"
	Stdout.Println("game.loadConfig(",filename,")")
	
	
	
	// create the main application log
	mainFile, err := os.Create("main.log")
	if (err != nil) {
		Stderr.Println(err.Error())
		Mainlog = Stdout
	} else {
		Mainlog = log.New(mainFile, "", flags)
	}	

	// b, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// err = json.Unmarshal(b, &territories)
	// if (err != nil) {
	// 	panic(err)
	// }
	
}

