package main

import (
	"fmt"
)

//闭包的本质
//     闭包是包含自由变量的代码块，这些变量不在这个代码块内或者任何全局上下文中定义，
//  而是在定义代码块的环境中定义。由于自由变量包含在代码块中，所以只要闭包还被使用，
//  那么这些自由变量以及它们引用的对象就不会被释放，要执行的代码为自由变量提供绑定的计算环境。

//     闭包的价值在于可以作为函数对象或者匿名函数，对于类型系统而言，这意味着不仅要表示数据还要表示代码。
//  支持闭包的多数语言都将函数作为第一级对象，就是说这些函数可以存储到变量中作为参数传递给其他函数，
//  最重要的是能够被函数动态创建和返回。

//Golang中的闭包同样也会引用到函数外的变量，闭包的实现确保只要闭包还被使用，
//那么被闭包引用的变量会一直存在。从形式上看，匿名函数都是闭包。

//闭包=函数+引用环境
//引用环境的定义：在程序执行中的某个点所有处于活跃状态的约束所组成的集合，
//             其中的约束指的是一个变量的名字和其所代表的对象之间的联系。

//====================函数作为参数=====================
// CalcFunc type
type CalcFunc func(x, y int) int

// AddFunc function
func AddFunc(x, y int) int {
	return x + y
}

// SubFunc function
func SubFunc(x, y int) int {
	return x - y
}

// OperationFunc function is params
func OperationFunc(x, y int, calcFunc CalcFunc) int {
	return calcFunc(x, y)
}

//================函数作为返回值==============================
//第一种写法
func add(x, y int) func() int {
	f := func() int {
		return x + y
	}
	return f
}

//第二中写法
func add2(x, y int) func() int {
	return func() int {
		return x + y
	}
}

//当函数返回多个匿名函数时建议采用第一种写法：

func calc(x, y int) (func(int) int, func() int) {
	f := func(z int) int {
		return (x + y) * z / 2
	}
	f2 := func() int {
		return 2 * (x + y)
	}
	return f, f2
}

// 在匿名函数定义的同时进行调用：花括号后跟参数列表表示函数调用
func safeHandler() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("some exception has happend:", err)
		}
	}()
}

func adds(n int) func(int) int {
	sum := n
	f := func(x int) int {
		i := 2
		sum += i * x
		return sum
	}
	return f
}

//   该例子中函数变量为f，自由变量为sum，同时f为sum提供绑定的计算环境，使得sum和f粘滞在了一起，它们组成的代码块就是闭包。
//add函数的返回值是一个闭包，而不仅仅是f函数的地址。在该闭包函数中，只有内部的匿名函数f才能访问局部变量i，
//而无法通过其他途径访问，因此闭包保证了i的安全性。
//   当我们分别用不同的参数(10, 20)注入add函数而得到不同的闭包函数变量时，得到的结果是隔离的，
//也就是说每次调用add函数后都将生成并保存一个新的局部变量sum。
// 名言：对象是附有行为的数据，而闭包是附有数据的行为

func main() {
	//===========1=========================
	sum := OperationFunc(1, 2, AddFunc)
	fmt.Println(sum) //3
	difference := OperationFunc(1, 2, SubFunc)
	fmt.Println(difference) // -1
	//============2========================
	f1 := add(1, 2)
	fmt.Println(f1())
	f2 := add2(1, 2)
	fmt.Println(f2())
	ff1, ff2 := calc(2, 3)
	n1 := ff1(10)
	n2 := ff1(20)
	n3 := ff2()
	fmt.Printf("n1 = %d, n2= %d, n3 = %d\n", n1, n2, n3)
	//n1 = 25, n2= 50, n3 = 10

	ffs := adds(10)
	nns := ffs(3) //16
	fmt.Println(nns)
	nns2 := ffs(4)
	fmt.Println(nns2) //16+2*4 = 24
	nns3 := ffs(2)
	fmt.Println(nns3) //24+2*2 = 28
	fmt.Println("==============")
	ffm := adds(20)
	mms := ffm(3)
	fmt.Println(mms)    //20 + 2*3 = 26
	fmt.Println(ffm(2)) //26 + 2*2 = 30
}
