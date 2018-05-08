package flaws

import "testing"
import "bytes"

func BenchmarkUnknown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.Write([]byte{1})
		_ = buf.Bytes()
	}
}

//cmd / compile，bytes：bootstrap数组导致bytes.Buffer始终被堆分配
