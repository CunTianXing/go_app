package main

import (
	"sync"
	"testing"
)

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

var pool = sync.Pool{
	New: func() interface{} {
		return &Book{}
	},
}

func BenchmarkNoPool(b *testing.B) {
	var book *Book

	for n := 0; n < b.N; n++ {
		book = &Book{
			Title:  "xingcuntian test",
			Author: "Gary",
			Pages:  100,
		}
	}

	_ = book
}

func BenchmarkPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		book := pool.Get().(*Book)
		book.Title = "xingcuntian test"
		book.Author = "Gary"
		book.Pages = 100
		pool.Put(book)
	}
}
