package main

import (
	"fmt"
	"reflect"
)

type I interface {
	M() string
}

type T struct {
	name string
}

func (t T) M() string {
	return t.name
}

func Hello(i I) {
	fmt.Printf("Hi, my name is %s\n", i.M())
}

//在函数Hello中，方法调用i.M()通过特定的方式实现：当方法是类型满足的接口的实现时，不同类型的方法可以被调用。

//Golang的重要特征：接口是隐式实现的，程序员不需要显示声明类型T实现了接口I。这项工作是由Go编译器自动完成的(永远不要让人去做机器应该做的事情)。
//   这种行为的优雅实现使得如下这种方式成为可能：定义一个接口，这个接口被已经写好的类型自动实现(不需要对之前已完成的类型做修改)。
//这种语言级的特性，可以在新增接口时，既有类型不修改，则自动实现了新的接口。这种多态的方式具有很好的灵活性。

type T1 struct {
	name string
}

func (t T1) M() string {
	return t.name
}

type T2 struct {
	name string
}

func (t T2) M() string {
	return t.name + "hhhhhhh"
}

func Hi(i I) {
	fmt.Printf("Hi, my name is %s\n", i.M())
}

//单个类型可以实现多个接口
type I1 interface {
	M1()
}
type I2 interface {
	M2()
}

type TT struct{}

func (TT) M1() {
	fmt.Println("TT.I1")
}
func (TT) M2() {
	fmt.Println("TT.I2")
}

func f1(i I1) { i.M1() }
func f2(i I2) { i.M2() }

type F interface {
	M()
}

type TF struct{}

func (TF) M() {} //类型实现了接口F

func main() {
	Hello(T{"xingcuntian"})
	fmt.Println("ok")
	//同一个接口可以被多种类型实现
	Hi(&T1{name: "xingcuntian"})
	Hi(T2{name: "xingcuntian"})
	//单个类型可以实现多个接口
	f1(TT{})
	f2(TT{})
	//除了一个或多个接口要求的方法，类型可以自由实现其他不同的方法。

	//接口I类型的变量可以保存任何实现了接口I的值
	var i I = T{"name"}
	fmt.Println(i)
	fmt.Printf("value %#v\n", i)

	var ii I = T1{}        // 接口类型变量i，可以赋值为T1，也可以被赋值为类型T2
	fmt.Printf("%T\n", ii) // 输出main.T1
	ii = T2{}
	fmt.Printf("%T\n", ii) // 输出main.T2
	_ = ii

	fmt.Println(reflect.TypeOf(i).PkgPath(), reflect.TypeOf(i).Name())
	fmt.Println(reflect.TypeOf(i).String())
	fmt.Printf("%T\n", i)
	//接口类型空值nil
	var f *TF      // 变量f必然是空值nil
	fmt.Println(f) // <nil>
	fmt.Printf("f type: %T\n", f)
	if f == nil {
		fmt.Println("f is nil") //  f is nil
	} else {
		fmt.Println("f is not nil")
	}

	var ti F = f
	fmt.Println(ti) //<nil>
	fmt.Printf("ti type: %T\n", ti)
	if ti == nil {
		fmt.Println("ti is nil")
	} else {
		fmt.Println("ti is not nil") //ti is not nil
	}

	//接口类型值的动态类型为*main.T，所有它不等于空值nil
}
