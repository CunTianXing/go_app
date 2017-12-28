package main

import (
	"fmt"
	"log"

	"github.com/xiam/exif"
)

func main() {
	data, err := exif.Read("image/5.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	for key, val := range data.Tags {
		fmt.Printf("%s = %s\n", key, val)
	}
}
