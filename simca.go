package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "net"
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

    //IOLoop()
    TCPLoop()

    fmt.Printf("Come back soon!\n")
}

func TCPLoop() {
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
        return
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Printf("Error accepting: %s\n", err)
            continue
        }

        fmt.Printf("New connection from %s\n", conn.RemoteAddr())
        go handleRequest(conn)
    }
}

func handleRequest(conn net.Conn) {
    var text, key, val, response_str string
    var response []byte

    IOLoop:
    for {
        request := make([]byte, 65536)
        n, err := conn.Read(request)
        if err != nil {
            fmt.Printf("Error reading: %s\n", err)
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, request)

        // TODO: Fix ugly bug with string(byte(0))
        text = string(request)
        r := strings.NewReplacer("\r", "", "\n", "", string(byte(0)), "")
        text = r.Replace(text)
        args := strings.Split(text, " ")

        switch args[0] {
            case "set":
                if len(args) < 3 {
                    response_str = fmt.Sprintf("set() expects 2 args, %d given\n", len(args) - 1)
                    break
                }
                key = args[1]
                val = strings.Join(args[2:], " ")

                set(key, val)
                // TODO: Check for errors on SET
                response_str = fmt.Sprintf("OK\n")
                break;

            case "get":
                if len(args) < 2 {
                    response_str = fmt.Sprintf("get() expects 1 arg, %d given\n", len(args) - 1)
                    break
                }
                key = args[1]
                val = get(key)
                response_str = fmt.Sprintf("Data at key '%s' is '%s'\n", key, val)
                break;

            case "len":
                response_str = fmt.Sprintf("Data length is %d\n", len(keys_hashmap))
                break

            case "dump":
                response_str = fmt.Sprintf("DUMP is not implemented for TCP yet\n")
                /*for key, elem := range keys_hashmap {
                    fmt.Printf("%s -> %s\n", key, (*elem).val)
                }*/
                break

            case "exit":
                break IOLoop

            default:
                response_str = fmt.Sprintf("Unknown command '%s'\n", args[0])
        }

        response = []byte(response_str)
        n, err = conn.Write(response)
        if err != nil {
            fmt.Printf("Error writing: %s\n", err)
            break
        }
    }

    n, err := conn.Write([]byte("Come back soon!\n"))
    if err != nil {
        fmt.Printf("Error writing after %d bytes: %s\n", n, err)
    }

    err = conn.Close()
    if err != nil {
        fmt.Printf("Error closing: %s\n", err)
    }
}

func IOLoop() {
    for true {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSuffix(text, "\n")

        var key, val string
        args := strings.Split(text, " ")

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

            case "dump":
                for key, elem := range keys_hashmap {
                    fmt.Printf("%s -> %s\n", key, (*elem).val)
                }
                break

            case "":
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