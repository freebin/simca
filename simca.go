package main

import (
    "fmt"
    "net"
    "strings"
)

func main() {
    fmt.Printf("Starting S.I.M.C.A.\n")
    res := initStorage()
    if !res {
        fmt.Printf("Failed to init storage!\n")
        return
    }

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
            // TODO: Fix "Error reading: EOF" or replace with binary protocol
            fmt.Printf("Error reading: %s\n", err)
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, request)

        // TODO: Fix ugly bug with string(byte(0)) or replace with binary protocol
        text = string(request)
        r := strings.NewReplacer("\r", "", "\n", "", string(byte(0)), "")
        text = r.Replace(text)
        args := strings.Split(text, " ")

        response_str = ""

        switch args[0] {
            case "set":
                if len(args) < 3 {
                    response_str = fmt.Sprintf("set() expects 2 args, %d given\n", len(args) - 1)
                    break
                }
                key = args[1]
                val = strings.Join(args[2:], " ")

                set(key, val)
                response_str = fmt.Sprintf("OK\n")
                break

            case "get":
                if len(args) < 2 {
                    response_str = fmt.Sprintf("get() expects 1 arg, %d given\n", len(args) - 1)
                    break
                }
                key = args[1]
                val, found := get(key)
                if found {
                    response_str = fmt.Sprintf("Data at key '%s' is '%s'\n", key, val)
                } else {
                    response_str = fmt.Sprintf("Data at key '%s' is not found\n", key)
                }
                break

            case "len":
                response_str = fmt.Sprintf("Data length is %d\n", countEntries())
                break

            case "flush":
                flushStorage()
                response_str = fmt.Sprintf("FLUSHED! I hope you were not kidding...\n")
                break

            case "dump":
                response_str = printFullList()
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
