package main

import (
	"fmt"
	"lex"
	"strings"
)

func main() {
	lex.NewJslLexer(strings.NewReader("foobar"))
	fmt.Printf("%s", "Done")
}
