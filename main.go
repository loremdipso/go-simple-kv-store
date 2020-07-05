package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// turn off debug messages
	log.SetOutput(ioutil.Discard)

	store, err := NewStore("data.db")
	if err != nil {
		log.Println("{}", err)
	} else {
		defer store.Close()
		store.Parse(os.Args[1:])
	}
}
