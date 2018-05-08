package main

import "fmt"

func main() {
	arr := [10]int{}
	for i := 0; i < 10; i++ {
		fmt.Print("Result of ", i, ":")
		go func() {
			arr[i] = i + i*i
			fmt.Println(arr[i])
		}()
	}
	fmt.Println("Done")
}
