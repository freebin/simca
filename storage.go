package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
)

type Elem struct {
	Key string
	Val string
}

const permFilePath = "/tmp/simca.gob"

var keysHashmap map[string]Elem

func initStorage() bool {
	permFile, err := os.OpenFile(permFilePath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Failed to read datafile: %s!\n", err)
		return false
	}
	defer permFile.Close()

	decoder := gob.NewDecoder(permFile)
	err = decoder.Decode(&keysHashmap)
	if err == io.EOF {
		keysHashmap = make(map[string]Elem)
	} else if err != nil {
		fmt.Printf("Failed to import datafile: %s!\n", err)
		return false
	}

	// TODO: Save to file from time to time

	return true
}

func storeToFile() bool {
	permFile, err := os.OpenFile(permFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Failed to open datafile: %s!\n", err)
		return false
	}
	defer permFile.Close()

	encoder := gob.NewEncoder(permFile)
	err = encoder.Encode(keysHashmap)
	if err != nil {
		fmt.Printf("Failed to write datafile: %s!\n", err)
	}

	return true
}

func set(key string, val string) {
	//var existingElem Elem = find(key)
	var new_elem Elem = Elem{key, val}
	keysHashmap[key] = new_elem

	storeToFile()
}

func get(key string) (string, bool) {
	var existingElem Elem = find(key)
	if existingElem.Key == "" {
		return "", false
	}

	return existingElem.Val, true
}

func find(key string) Elem {
	//TODO: Better use link
	return keysHashmap[key]
}

func countEntries() int {
	return len(keysHashmap)
}

func flushStorage() {
	keysHashmap = make(map[string]Elem)
}

func printFullList() string {
	var dump string
	for key, elem := range keysHashmap {
		dump += fmt.Sprintf("%s -> %s\n", key, elem.Val)
	}
	return dump
}