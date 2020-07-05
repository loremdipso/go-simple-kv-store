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
	fmt.Println("Closing database...")
	// log.Println("Closing database...")
	// NOTE: this takes up some time, so maybe don't do it always?
	// Or, maybe do it synchronously?
	store.bdb.Close()
}

// Parse some key-store queries
func (store *KVStore) Parse(args []string) {
	if len(args) > 0 {
		if args[0] == "get" {
			store.tryGetMultiple(args[1:])
		} else if args[0] == "set" {
			store.trySetMultiple(args[1:])
		} else if args[0] == "setblobs" {
			store.trySetMultipleBlobs(args[1:])
		} else if args[0] == "clean" {
			store.clean()
		} else {
			fmt.Printf("ERROR: %s not a recognized command\n", args[0])
		}
	}
}

func (store *KVStore) clean() {
	store.bdb.Shrink()
}

func (store *KVStore) tryGetMultiple(ids []string) {
	for _, id := range ids {
		value, err := store.getSingle(id)
		if err != nil {
			fmt.Printf("Error for %s: %s\n", id, err)
		} else {
			fmt.Printf("Result for %s: %s\n", id, value)
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
				fmt.Printf("Setting value for %s\n", id)
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
	err := store.bdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	return err
}
