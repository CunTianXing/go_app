package main

import (
    "fmt"
    "os"
    "os/signal"
    "time"
)

func main() {
    ticker := time.NewTicker(3*time.Second)
    defer ticker.Stop()
    c := make(chan os.Signal,1)
    signal.Notify(c,os.Interrupt)
    for {
        select {
           case now := <-ticker.C:
               fmt.Printf("tick %s\n", now.UTC().Format("20060102-150405.000000000"))
           case <-c:
              fmt.Println("Received C-c - shutting down")
              return
        }
    }
}
