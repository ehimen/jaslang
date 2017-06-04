package main

import (
	"bytes"
	"strings"

	"fmt"

	"github.com/ehimen/jaslang/run"
)

func main() {
	code := strings.NewReader(`println(1 + 2);`)
	input := strings.NewReader("")
	output := bytes.NewBufferString("")

	if err := run.Interpret(code, input, output); err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Println(output.String())
}
