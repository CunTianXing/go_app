package main

import "fmt"

//*new(T)等于nil 如果T型的零值用nil表示
func main() {
	fmt.Println(*new(*int) == nil)
	fmt.Println(*new([]int) == nil)
	fmt.Println(*new(map[string]int) == nil)
	fmt.Println(*new(chan bool) == nil)
	fmt.Println(*new(func()) == nil)
	fmt.Println(*new(interface{}) == nil)
}
