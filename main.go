package main

import (
	"log"
	"os"
)

func main() {
	store, err := NewStore("data.db")
	if err != nil {
		log.Println("{}", err)
	}
	defer store.Close()
	store.Parse(os.Args[1:])
}
