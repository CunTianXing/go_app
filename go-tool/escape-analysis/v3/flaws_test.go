package flaws

import "testing"

func BenchmarkLiteralFunctions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var y1 int
		foo(&y1, 42) //GOOD: y1 does not escape

		var y2 int
		func(p *int, x int) {
			*p = x
		}(&y2, 42) // GOOD: y2 does not escape

		var y3 int
		p := foo
		p(&y3, 42) //BAD: Cause of y3 escape
	}
}

func foo(p *int, x int) {
	*p = x
}

//foo函数被分配给一个名为变量的变量p。通过p变量，该foo函数与y3变量共享执行。这个函数调用是通过p变量的间接完成的
// y3分配在堆上
