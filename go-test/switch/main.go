package main

import (
	"fmt"
)

func main() {
	switchOne(3)
	switchTwo()
	multipleSwitch("a")
	multipleSwitch("A")
	expressionSwitch(9)
	expressionSwitch(20)
	funcSwitch()
}

func switchOne(finger int) {
	switch finger {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	case 4:
		fmt.Println("four")
	case 5:
		fmt.Println("five")
	}
}

func switchTwo() {
	switch finger := 4; finger {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	default:
		fmt.Println("default")
	}
}

func multipleSwitch(pointor string) {
	switch pointor {
	case "a", "b", "c", "d":
		fmt.Println("abcd")
	default:
		fmt.Println("no pointer")
	}
}

func expressionSwitch(num int) {
	switch {
	case num > 0 && num < 10:
		fmt.Println("[1:10]")
	case num < 100 && num > 10:
		fmt.Println("[11:100]")
	}
}

func number() int {
	num := 20 * 2
	return num
}

func funcSwitch() {
	switch num := number(); {
	case num <= 10 && num >= 0:
		fmt.Println("num is 0 -- 10")
	case num > 10 && num < 100:
		fmt.Println("num is 10--100")
	}

}
