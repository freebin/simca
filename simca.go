package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    fmt.Printf("Starting S.I.M.C.A.\n")

IOLoop:
    for true {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("> ")
        text, _ := reader.ReadString('\n')
        text = strings.TrimSuffix(text, "\n")

        switch text {
            case "set":
                break;

            case "get":
                break;

            case "exit":
                break IOLoop

            default:
                fmt.Printf("Unknown command '%s'\n", text)
        }
    }

    fmt.Printf("Come back soon!\n")
}
