package main

import (
    "fmt"
    "net"
    "strings"
    "time"
)

type InternalStats struct {
    StartTime uint
    NumConnections uint
    NumGetRequests uint
    NumGetHits uint
    NumSetRequests uint
    NumSetHits uint
    NumKeysFromStorage uint
}

var stats InternalStats

func main() {
    stats.StartTime = uint(time.Now().Unix())

    fmt.Printf("Starting S.I.M.C.A.\n")
    res := initStorage()
    if !res {
        fmt.Printf("Failed to init storage!\n")
        return
    }

    stats.NumKeysFromStorage = uint(countEntries())

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

    stats.NumConnections++

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

                _, found := set(key, val)
                if found {
                    stats.NumSetHits++
                }
                response_str = fmt.Sprintf("OK\n")
                stats.NumSetRequests++
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
                    stats.NumGetHits++
                } else {
                    response_str = fmt.Sprintf("Data at key '%s' is not found\n", key)
                }
                stats.NumGetRequests++
                break

            case "stats":
                response_str = ""
                response_str += fmt.Sprintf("StartTime: %d\n", stats.StartTime)
                response_str += fmt.Sprintf("NumKeysFromStorage: %d\n", stats.NumKeysFromStorage)
                response_str += fmt.Sprintf("DataLength: %d\n", countEntries())
                response_str += fmt.Sprintf("NumConnections: %d\n", stats.NumConnections)
                response_str += fmt.Sprintf("NumGetRequests: %d\n", stats.NumGetRequests)
                response_str += fmt.Sprintf("NumGetHits: %d\n", stats.NumGetHits)
                response_str += fmt.Sprintf("NumSetRequests: %d\n", stats.NumSetRequests)
                response_str += fmt.Sprintf("NumSetHits: %d\n", stats.NumSetHits)
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
