package main

import "fmt"

//It Is Legal To Range Over Nil Channels, Maps, Slices, And Array Pointers
// 迭代nil map和slice的循环次数为零。
// 通过迭代 nil指针数组的循环次数是其相应数组类型的长度。
//（但是，如果相应数组类型的长度不为零，并且第二次迭代既不会被忽略也不会被忽略，那么迭代在运行时会出现混乱。）

func main() {
	for range []int(nil) { // 0
		fmt.Println("Hello")
	}

	for range map[string]string(nil) { //0
		fmt.Println("world")
	}

	for i := range (*[5]int)(nil) { //0,1,2,3,4
		fmt.Println(i)
	}

	for range chan bool(nil) { // block here  all goroutines are asleep - deadlock!
		fmt.Println("Bye")
	}
}
