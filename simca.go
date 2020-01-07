package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Elem struct {
    key string
    val string
    prev *Elem
    next *Elem
}

var first_elem_pointer *Elem = nil
var last_elem_pointer *Elem = nil
var keys_hashmap map[string]*Elem

func main() {
    fmt.Printf("Starting S.I.M.C.A.\n")
    keys_hashmap = make(map[string]*Elem)

    IOLoop()

    fmt.Printf("Come back soon!\n")
}

func IOLoop() {
    for true {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSuffix(text, "\n")

        var key, val string
        args := strings.Split(text, " ")

        fmt.Printf("CMD is '%s'\n", args[0])

        switch args[0] {
            case "set":
                if len(args) < 3 {
                    fmt.Printf("set() expects 2 args, %d given\n", len(args) - 1)
                    break
                }
                key = args[1]
                val = strings.Join(args[2:], " ")

                set(key, val)
                break;

            case "get":
                if len(args) < 2 {
                    fmt.Printf("get() expects 1 arg, %d given\n", len(args) - 1)
                    break
                }
                key = args[1]
                val = get(key)
                fmt.Printf("Data at key '%s' is '%s'\n", key, val)
                break;

            case "len":
                fmt.Printf("Data length is %d\n", len(keys_hashmap))
                break

            case "exit":
                return

            default:
                fmt.Printf("Unknown command '%s'\n", args[0])
        }
    }
}

func set(key string, val string) {
    var existing_elem_pointer *Elem = find(key)
    if existing_elem_pointer == nil {
        var new_elem Elem = Elem{key, val, nil, nil}
        if last_elem_pointer == nil {
            first_elem_pointer = &new_elem
            last_elem_pointer = &new_elem
        } else {
            (*last_elem_pointer).next = &new_elem
            new_elem.prev = last_elem_pointer
            last_elem_pointer = &new_elem
        }
        keys_hashmap[key] = &new_elem
    } else {
        (*existing_elem_pointer).val = val
    }
}

func get(key string) string {
    // TODO: Distinguish NOT FOUND vs EMPTY STRING
    var existing_elem_pointer *Elem = find(key)
    if existing_elem_pointer == nil {
        return ""
    } else {
        return (*existing_elem_pointer).val
    }
}

func find(key string) *Elem {
    return keys_hashmap[key]
}