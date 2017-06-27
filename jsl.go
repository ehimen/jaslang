package main

import (
	"bytes"
	"strings"

	"fmt"

	"flag"

	"os"

	"bufio"

	"log"

	"github.com/ehimen/jaslang/run"
)

func main() {
	flag.Parse()

	file := flag.Arg(0)

	if len(file) == 0 {
		log.Fatal("Must specify jaslang source file (.jsl)")
	}

	if f, err := os.Open(file); err != nil {
		log.Fatal(err)
	} else {
		code := bufio.NewReader(f)

		input := strings.NewReader("")
		output := bytes.NewBufferString("")
		outputError := bytes.NewBufferString("")

		if run.Interpret(code, input, output, outputError) {
			log.Fatal(fmt.Sprintf("%v\n", outputError.String()))
		} else {
			fmt.Println(output.String())
		}
	}
}
