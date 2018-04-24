package main

import "fmt"

//Call Methods Through nil Arguments Will Not Panic
type Slice []bool

func (s Slice) Length() int {
	return len(s)
}

func (s Slice) Modify(i int, x bool) {
	fmt.Println(11111)
	s[i] = x // panic if s is nil
}

func (p *Slice) DoNothing() {

}

func (p *Slice) Append(x bool) {
	fmt.Println(22222)
	*p = append(*p, x) // panic if p is nil
}

//在Go中，使用nil作为方法接收器参数绝不是导致恐慌的直接原因。
func main() {
	//以下选择器不会导致恐慌。
	_ = ((Slice)(nil)).Length
	_ = ((Slice)(nil)).Modify
	_ = ((*Slice)(nil)).DoNothing
	_ = ((*Slice)(nil)).Append
	//以下两行也不会恐慌。
	_ = ((Slice)(nil)).Length()
	((*Slice)(nil)).DoNothing()
	//以下两行将会出现恐慌。 但是在调用方法时不会引发恐慌。 它将在方法体中触发。
	((Slice)(nil)).Modify(0, true)
	((*Slice)(nil)).Append(true)

}
