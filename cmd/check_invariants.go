package main

import (
    "fmt"
    "log"

    "indexer/invariants"
)

func main() {
    if err := invariants.CheckAll(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("All invariants passed ")
}
