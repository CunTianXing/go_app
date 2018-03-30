package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type line struct {
	id         int
	restOfLine string
}

func reader(fname string, out chan<- *line) {
	defer close(out)
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("open: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	header := true
	for scanner.Scan() {
		var l line
		columns := strings.SplitN(scanner.Text(), ",", 2)
		if header {
			header = false
			continue
		}
		id, err := strconv.Atoi(columns[0])
		if err != nil {
			log.Fatalf("ParseInt:%v", err)
		}
		l.id = id
		l.restOfLine = columns[1]
		out <- &l
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func joiner(metadata, setIDs <-chan *line, out chan<- *line) {
	defer close(out)
	si := &line{}
	var m []*line
	for md := range metadata {
		sep := ","
		fmt.Printf("md: %#v\n", md)
		if si.id == md.id {
			fmt.Printf("eq md: %#v\n", md)
			md.restOfLine += sep + si.restOfLine
			sep = " "
		} else if si.id > 0 {
			m = append(m, si)
		}
		for _, v := range m {
			if v.id == md.id {
				md.restOfLine += sep + v.restOfLine
			}
		}
		for si = range setIDs {
			fmt.Printf("si: %#v\n", si)
			if si.id == md.id {
				fmt.Printf("eq si: %#v\n", si)
				md.restOfLine += sep + si.restOfLine
				sep = " "
			} else if si.id > md.id {
				fmt.Println("break")
				break
			}
		}
		fmt.Printf("si end md: %#v\n", md)
		fmt.Println("=====================")
		out <- md
	}
}

func main() {
	metadataChan := make(chan *line)
	go reader("metadata.csv", metadataChan)
	fmt.Println("1")
	strengthSetChan := make(chan *line)
	go reader("strength_sets.csv", strengthSetChan)
	fmt.Println("2")
	augmentedLinesChan := make(chan *line)
	go joiner(metadataChan, strengthSetChan, augmentedLinesChan)
	fmt.Println("id,user_id,machine_id,circle_id,timestamp,strength_set_ids")
	for l := range augmentedLinesChan {
		fmt.Printf("%v, %v\n", l.id, l.restOfLine)
		fmt.Println("======end==================")
	}
}
