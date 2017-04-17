package main

import (
	"fmt"

	"github.com/ehimen/jaslang/parse"
)

func main() {
	var foo parse.ParentNode

	foo = &parse.FunctionCall{}

	fmt.Printf("%v", foo)
}
