package flaws

import (
	"net/http"
	"testing"
)

func BenchmarkHandler(b *testing.B) {
	h := func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}

	route := wrapHandler(h)

	for i := 0; i < b.N; i++ {
		var r http.Request
		route(nil, &r) //BAD: Cause of r escape
	}
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func wrapHandler(h Handler) Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		//fmt.Println("testing")
		return h(w, r)
	}
	return f
}
