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
}

var data []Elem

func main() {
    fmt.Printf("Starting S.I.M.C.A.\n")

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
                fmt.Printf("Data length is %d\n", len(data))
                break

            case "exit":
                return

            default:
                fmt.Printf("Unknown command '%s'\n", args[0])
        }
    }
}

func set(key string, val string) {
    var exists_on_index int = getIndexByKey(key)
    if exists_on_index == -1 {
        data = append(data, Elem{key, val})
    } else {
        data[exists_on_index] = Elem{key, val}
    }
}

func get(key string) string {
    // TODO: Distinguish NOT FOUND vs EMPTY STRING
    var exists_on_index int = getIndexByKey(key)
    if exists_on_index == -1 {
        return ""
    } else {
        return data[exists_on_index].val
    }
}

func getIndexByKey(key string) int {
    for i := 0; i < len(data); i++ {
        if data[i].key == key {
            return i;
        }
    }

    return -1;
}