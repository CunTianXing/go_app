package main

import "fmt"

func main() {
	//Two nil Values Of The Same Type May Be Not Comparable
	//在Go中，map。 切片和函数类型不支持比较。 因此，比较使用任何类型的不可比类型指定的两个零标识符是非法的
	// var _ = ([]int)(nil) == ([]int)(nil)
	// var _ = (map[string]int)(nil) == (map[string]int)(nil)
	// var _ = (func())(nil) == (func())(nil)
	//但是，上述无法比较的类型的任何值都可以与裸nil标识符进行比较。
	var _ = ([]int)(nil) == nil
	if ([]int)(nil) == nil {
		fmt.Println(true) // ok
	} else {
		fmt.Println(false)
	}
	var _ = (map[string]int)(nil) == nil
	if (map[string]int)(nil) == nil {
		fmt.Println(true) // ok
	} else {
		fmt.Println(false)
	}
	var _ = (func())(nil) == nil
	if (func())(nil) == nil {
		fmt.Println(true) // ok
	} else {
		fmt.Println(false)
	}
	//Two nil Values May Be Not Equal
	// 如果两个比较零值之一是一个接口值，另一个不是，假设它们是可比较的，那么比较结果总是假的。
	// 原因是在进行比较之前，非接口值将被转换为接口值的类型。 转换的接口值具有具体的动态类型，但其他接口值没有。 这就是为什么比较结果总是错误的原因。
	fmt.Println((interface{})(nil) == (*int)(nil)) // false
	// Retrieving Elements From Nil Maps Will Not Panic
	// Retrieving element from a nil map value will always return a zero element value.

	fmt.Println((map[string]int)(nil)["key"]) // 0
	fmt.Println((map[int]bool)(nil)[123])     // false
	fmt.Println((map[int]*int64)(nil)[1234])  //nil

}
