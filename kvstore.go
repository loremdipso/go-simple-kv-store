package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tidwall/buntdb"
)

type KVStore struct {
	bdb *buntdb.DB
}

type Blob struct {
	Id json.Number
}

// New KVStore at filename
func NewStore(filename string) (*KVStore, error) {
	buntDB, err := buntdb.Open(filename)
	if err != nil {
		return nil, err
	}

	return &KVStore{bdb: buntDB}, nil
}

// Close the store safely
func (store *KVStore) Close() {
	log.Println("Closing database...")
	// log.Println("Closing database...")
	// NOTE: this takes up some time, so maybe don't do it always?
	// Or, maybe do it synchronously?
	store.bdb.Close()
}

// Parse some key-store queries
func (store *KVStore) Parse(args []string) {
	if len(args) > 0 {
		if args[0] == "help" || args[0] == "-help" || args[0] == "--help" {
			fmt.Println("help")
			fmt.Println("get")
			fmt.Println("set: pass in pairs of key/value strings")
			fmt.Println("setblobs: sets a json to the key defined by the \"Id\" property in that json")
			fmt.Println("keys: get all keys")
			fmt.Println("dump: print all keys and values")
		} else if args[0] == "get" {
			store.tryGetMultiple(args[1:])
		} else if args[0] == "set" {
			store.trySetMultiple(args[1:])
		} else if args[0] == "setblobs" {
			store.trySetMultipleBlobs(args[1:])
		} else if args[0] == "clean" {
			store.clean()
		} else if args[0] == "dump" {
			store.dump()
		} else if args[0] == "keys" {
			store.dumpKeys()
		} else {
			fmt.Printf("ERROR: %s not a recognized command\n", args[0])
		}
	}
}

func (store *KVStore) clean() {
	store.bdb.Shrink()
}

func (store *KVStore) dumpKeys() {
	store.bdb.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys("*", func(key, _ string) bool {
			fmt.Println(key)
			return true
		})
		return nil
	})
}

func (store *KVStore) dump() {
	store.bdb.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys("*", func(key, value string) bool {
			fmt.Printf("%s: %s\n", key, value)
			return true
		})
		return nil
	})
}

func (store *KVStore) tryGetMultiple(ids []string) {
	for i, id := range ids {
		value, err := store.getSingle(id)
		if err != nil {
			log.Printf("Error for %s: %s\n", id, err)
		} else {
			// log.Printf("Result for %s: %s\n", id, value)
			if i > 0 {
				fmt.Printf("\n\n=========\n\n")
			}
			fmt.Println(value)
		}
	}
}

func (store *KVStore) trySetMultiple(pairs []string) {
	for i := 0; i < len(pairs)-1; i += 2 {
		key := pairs[i]
		value := pairs[i+1]
		store.setSingle(key, value)
	}
}

func (store *KVStore) trySetMultipleBlobs(blobs []string) {
	for _, blobstring := range blobs {
		var blob Blob
		err := json.Unmarshal([]byte(blobstring), &blob)
		if err != nil {
			log.Printf("ERROR deserializing %s\n", blobstring)
		} else {
			id := string(blob.Id)
			if len(id) == 0 {
				log.Printf("ERROR no id in : %s\n", blobstring)
			} else {
				store.setSingle(id, blobstring)
			}
		}
	}
}

func (store *KVStore) getSingle(key string) (string, error) {
	var value string
	err := store.bdb.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	return value, err
}

func (store *KVStore) setSingle(key, value string) error {
	log.Printf("Setting %s\n", key)
	err := store.bdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	return err
}
