package main

import (
	"bytes"
	"strings"

	"fmt"

	"github.com/ehimen/jaslang/run"
)

func main() {
	code := strings.NewReader(`println("Hello world!");`)
	input := strings.NewReader("")
	output := bytes.NewBufferString("")

	if err := run.Interpret(code, input, output); err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Println(output.String())
}
