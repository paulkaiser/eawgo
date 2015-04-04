// test the dice rolling package
package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"fmt"
	"strconv"
)

func main() {

	logFlags := log.Ldate | log.Lmicroseconds | log.Lshortfile
	log.SetFlags(logFlags)
	log.Println("logger prefix:", log.Prefix)
	log.Println("Start controller-test")
	
	parseEnvVars()
	
	bx, by := getBoardSize()
	log.Println("bx, by:", bx, by)
}


func parseEnvVars () map[string]string {
	
	// get all env vars that begin with EAW
	// store as key=value pairs in the returned map.
	pat := "^EAW_.+="
	varray := os.Environ()
	
	vars := make(map[string]string)
	
	for i := 0; i < len(varray); i++ {
		ev := varray[i]
		var m bool
		var err error
		m,err = regexp.MatchString(pat, ev)
		if (err != nil) {
			panic(fmt.Sprint("regexp.MatchString returned error",err,"for pat:",pat,"ev:",ev))
		}
		// log.Println("pat,ev,m,err",pat,ev,m,err)
		if (m) {
			subs := strings.SplitN(ev,"=", 2)
			vars[subs[0]] = subs[1]
			// log.Println("subs",subs)
		}
	}
	log.Println("vars",vars)
	return vars
	
}

func getBoardSize() (x int, y int) {
	
	var bx int
	var by int
	var err error
	var val string
	
	val = os.Getenv("EAW_BOARD_X")
	if val != "" {
		bx, err = strconv.Atoi(val)
		if (err != nil) {
			panic(err)
		}
	} else {
		bx = 2 //default
	}

	val = os.Getenv("EAW_BOARD_Y")
	if val != "" {
		by, err = strconv.Atoi(val)
		if (err != nil) {
			panic(err)
		}
	} else {
		by = 2 //default
	}
	
	return bx, by
}
