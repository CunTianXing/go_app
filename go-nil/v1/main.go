package main

import (
	"fmt"
)

//在Go中，nil可以表示以下类型的零值：
// pointer types (including type-unsafe ones).
// map types.
// slice types.
// function types.
// channel types.
// interface types.
func main() {
	//编译器必须有足够的信息来推导出一个nil值的类型。
	ok := (*struct{})(nil)
	fmt.Printf("ok is %#v\n", ok)
	_ = []int(nil)
	_ = map[int]bool(nil)
	_ = chan string(nil)
	_ = (func())(nil)
	_ = interface{}(nil)

	var _ *struct{} = nil
	var _ []int = nil
	var _ map[int]bool = nil
	var _ chan string = nil
	var _ func() = nil
	var _ interface{} = nil
	var i interface{}
	fmt.Println(i)
	//nil Is Not A Keyword In Go

	nil := 123
	fmt.Println(nil)
}
