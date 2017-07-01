package main

import (
	"bytes"
	"strings"

	"fmt"

	"flag"

	"os"

	"log"

	"io"

	"bufio"

	"encoding/json"

	"github.com/ehimen/jaslang/lex"
	"github.com/ehimen/jaslang/parse"
	"github.com/ehimen/jaslang/run"
)

func main() {
	ast := flag.Bool("ast", false, "Prints the parsed AST as JSON. Does not execute code")

	flag.Parse()

	file := flag.Arg(0)

	var input io.RuneReader

	if len(file) == 0 {
		input = bufio.NewReader(os.Stdin)
	} else {
		if f, err := os.Open(file); err != nil {
			log.Fatal(err)
		} else {
			input = bufio.NewReader(f)
		}
	}

	if *ast {
		printAst(input)
	} else {
		execute(input)
	}
}

func execute(file io.RuneReader) {
	input := strings.NewReader("")
	output := bytes.NewBufferString("")
	outputError := bytes.NewBufferString("")

	if run.Interpret(file, input, output, outputError) {
		fail(outputError.String())
	} else {
		fmt.Println(output.String())
	}
}

func printAst(file io.RuneReader) {
	parser := parse.NewParser(lex.NewJslLexer(file))

	if ast, err := parser.Parse(); err != nil {
		fail(err.Error())
	} else {
		if astJson, err := json.MarshalIndent(ast, "", "    "); err != nil {
			fail(err.Error())
		} else {

			fmt.Println(string(astJson))
		}
	}
}

func fail(msg string) {
	log.Fatal(fmt.Sprintf("%s\n", msg))
}
