package main

import (
	"bytes"

	"github.com/z7zmey/php-parser/php7"
	"github.com/z7zmey/php-parser/visitor"
)

func main() {
	src := bytes.NewBufferString(`<? echo "Hello world";`)
	nodes, comments, positions := php7.Parse(src, "example.php")

	visitor := visitor.Dumper{
		Indent:    "",
		Comments:  comments,
		Positions: positions,
	}
	nodes.Walk(visitor)
}
