package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

func main(){ 
    rand.Seed(time.Now().UnixNano())

    const N = 5
    var values [N]int32

    var wg sync.WaitGroup
    wg.Add(N)
     for i := 0; i < N; i++ {
        i := i
         go func(){
            values[i] = 50 + rand.Int31n(50)
            fmt.Println("Done:", i)
            wg.Done() // wg.Add(-1)
        }()
    }
    wg.Wait()
    fmt.Println("values: ",values)
}
