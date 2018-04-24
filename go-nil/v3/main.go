package main

import "fmt"

// 在Go中，只有一个值可以隐式转换为另一个类型时，才能比较两个不同类型的两个值。 具体而言，有三种情况可以比较两种不同可比较的两种值：
// 两个值中的一个的类型是另一个的基础类型。
// 两个值之一的类型实现另一个值的类型（必须是接口类型）。
// 两个值中的一个的类型是定向通道类型，另一个是双向通道类型，并且这两种类型具有相同的元素类型。
// 零值不是上述规则的例外。

// IntPtr  base *int
type IntPtr *int

func main() {
	// The underlying of type IntPtr is *int.
	var _ = IntPtr(nil) == (*int)(nil)
	if IntPtr(nil) == (*int)(nil) {
		fmt.Println(true) //ok
	}
	// Every type in Go implements interface{} type.
	var _ = (interface{})(nil) == (*int)(nil)
	if (interface{})(nil) == (*int)(nil) {
		fmt.Println(true)
	} else {
		fmt.Println(false) //ok
	}

	var _ = (chan int)(nil) == (chan<- int)(nil)
	var _ = (chan int)(nil) == (<-chan int)(nil)
	if (chan int)(nil) == (chan<- int)(nil) {
		fmt.Println(true) //ok
	} else {
		fmt.Println(false)
	}
}
