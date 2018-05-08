package main

import (
    "fmt"
    "runtime"
    "runtime/debug"
    "sync"
)

func main(){
    newFunc := func() interface{} {
        return make([]byte,32)
    }

    pool := sync.Pool{New: newFunc}

    v1 := pool.Get().([]byte)
    fmt.Printf("v1 : %v\n", v1)

    v1[0] = 1
    fmt.Printf("modified v1 : %v\n",v1)

    pool.Put(v1)
    v2 := pool.Get().([]byte)
    fmt.Printf("v2: %v\n", v2)
    pool.Put(v2)

    //after GC
    debug.SetGCPercent(100)

    runtime.GC()

    v3 := pool.Get().([]byte)
    fmt.Printf("v3: %v\n", v3)
}
