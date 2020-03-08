package main

import "fmt"

type Elem struct {
	key string
	val string
}

var keysHashmap map[string]*Elem

func initStorage() bool {
	keysHashmap = make(map[string]*Elem)
	return true
}

func set(key string, val string) {
	var existingElemPointer *Elem = find(key)
	if existingElemPointer == nil {
		var new_elem Elem = Elem{key, val}
		keysHashmap[key] = &new_elem
	} else {
		(*existingElemPointer).val = val
	}
}

func get(key string) (string, bool) {
	var existing_elem_pointer *Elem = find(key)
	if existing_elem_pointer == nil {
		return "", false
	}

	return (*existing_elem_pointer).val, true
}

func find(key string) *Elem {
	return keysHashmap[key]
}

func countEntries() int {
	return len(keysHashmap)
}

func flushStorage() {
	keysHashmap = make(map[string]*Elem)
}

func printFullList() string {
	var dump string
	for key, elem := range keysHashmap {
		dump += fmt.Sprintf("%s -> %s\n", key, (*elem).val)
	}
	return dump
}