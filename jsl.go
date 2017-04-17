package main

import (
	"fmt"

	"github.com/ehimen/jaslang/parse"
)

func main() {

	var foo parse.Node

	foo = &parse.Statement{}

	if f, ok := foo.(parse.ContainsChildren); ok {
		fmt.Printf("Foo: %T\n", f)
		fmt.Printf("Is parent node, %+v\n", foo)
	} else {
		fmt.Printf("Is not parent node, %T\n", foo)
	}
}