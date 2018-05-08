package flaws

import "testing"

func BenchmarkAssignmentIndirect(b *testing.B) {
	type X struct {
		p *int
	}

	for i := 0; i < b.N; i++ {
		var i1 int
		x1 := &X{
			p: &i1,
		}

		_ = x1

		var i2 int
		x2 := &X{}
		(*x2).p = &i2
	}
}
