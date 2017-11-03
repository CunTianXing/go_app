package main

import (
	"log"
	"os"
	"runtime/pprof"
)

func bigBytes() *[]byte {
	s := make([]byte, 100000000)
	return &s
}

func main() {
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()
	for i := 0; i < 10; i++ {
		s := bigBytes()
		if s == nil {
			log.Println("no")
		}
	}
}
