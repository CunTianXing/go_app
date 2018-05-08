package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	fd, err := strconv.Atoi(os.Getenv("TSAROUTFD"))
	if err != nil {
		fmt.Printf("Coudln't parse TSAROUTFD: %s\n", err)
		os.Exit(3)
	}

	outPipe := os.NewFile(uintptr(fd), "outPipe")

	scanner := bufio.NewScanner(outPipe)

	for scanner.Scan() {
		fmt.Println(time.Now().Truncate(time.Second), "-->", scanner.Text())
	}
}

//TSAROUTFD=1 go run main.go
