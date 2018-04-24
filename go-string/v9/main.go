package main

import (
    "fmt"
    "time"
)

func main(){
    bs := make([]byte, 1<<26)
    s0 := string(bs)
    s1 := string(bs)
    s2 := s1
    // s0, s1 and s2 are three equal strings.
	// The underlying bytes of s0 is a copy of bs.
	// The underlying bytes of s1 is also a copy of bs.
	// The underlying bytes of s0 and s1 are two different copies of bs.
	// s2 shares the same underlying bytes with s1.
    
    startTime :=time.Now()
    _ = s0 == s1
    duration := time.Now().Sub(startTime)
    fmt.Println("duration for (s0 === s1):",duration)

    startTime =time.Now()
    _ = s1 == s2
    duration = time.Now().Sub(startTime)
    fmt.Println("duration for (s1 === s2):",duration)
}
//duration for (s0 === s1): 12.81749ms
//duration for (s1 === s2): 1.219µs
