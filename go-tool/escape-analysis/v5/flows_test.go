package flaws

import "testing"

func BenchmarkSliceMapAssignment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[int]*int)
		var x1 int
		m[0] = &x1 // BAD: cause of x1 escape

		s := make([]*int, 1)
		var x2 int
		s[0] = &x2 // BAD: cause of x2 escape
	}
}

//“切片和map分配”缺陷与在切片或map内共享值时发生的分配有关
// 在7行上map，该map存储类型值的地址int。
// 然后在第8行，int在第9行的map中创建并共享一个类型的值，键值为0.
// 第11行的地址切片也会发生同样的情况。
// 创建切片后，类型int值为在索引0内共享。
