package run_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ehimen/jaslang/run"
)

// Relative to GOPATH
var jsltPath string = "src/github.com/ehimen/jaslang/jslt"

type codeBuffer interface {
	io.Seeker
	io.Reader
	io.RuneReader
}

type testCase struct {
	name   string
	code   codeBuffer
	output io.Reader
	input  io.Reader
	error  io.Reader
}

func (t *testCase) isValid() bool {
	return len(t.name) > 0 && t.code != nil && t.output != nil
}

func TestJslt(t *testing.T) {
	if tests, loaded := loadTests(t); !loaded {
		return
	} else {
		for _, test := range tests {
			runTest(t, test)
		}
	}
}

func runTest(t *testing.T, test testCase) {

	expected, _ := ioutil.ReadAll(test.output)

	actual := bytes.NewBufferString("")

	err := run.Interpret(test.code, test.input, actual)

	if actual.String() != string(expected) {
		test.code.Seek(0, io.SeekStart)
		code, _ := ioutil.ReadAll(test.code)
		t.Errorf(
			"\"%s\" failed!\nExpected output:\n%s\nActual output:\n%s\nErrors:\n%s\nCode:\n%s\n",
			test.name,
			expected,
			actual.String(),
			err,
			code,
		)
	}
}

func loadTests(t *testing.T) ([]testCase, bool) {
	goPath, exists := os.LookupEnv("GOPATH")

	if !exists {
		t.Log("Skipping JSLT tests as GOPATH is not defined")
		return nil, false
	}

	tests := []testCase{}

	for _, path := range filepath.SplitList(goPath) {
		path = filepath.Join(path, jsltPath)

		t.Log("Looking for jslt files in " + path)

		// TODO: potential duplication of file if identical
		// TODO: paths appear in GOPATH
		if files, err := filepath.Glob(filepath.Join(path, "*.jslt")); err != nil {
			return nil, false
		} else {
			for _, file := range files {
				if testsInFile, parseErr := parseFile(file); parseErr != nil {
					t.Errorf("Cannot read test file %s: %s", file, parseErr)
				} else {
					for _, test := range testsInFile {
						tests = append(tests, test)
					}
				}
			}
		}
	}

	return tests, true
}

func parseFile(path string) ([]testCase, error) {
	tests := []testCase{}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	fileReader := bufio.NewReader(file)

	prefixTest := "<!"
	prefixInput := "<<<INPUT"
	prefixOutput := "<<<OUTPUT"
	prefixCode := "<<<CODE"
	prefixError := "<<<ERROR"

	test := testCase{}
	content := ""
	current := ""

	closeSection := func(next string) error {
		if len(current) > 0 {
			// TODO: is it okay to ignore all trailing newlines?
			content = strings.Trim(content, "\n")
			switch current {
			case prefixInput:
				test.input = strings.NewReader(content)
			case prefixOutput:
				test.output = strings.NewReader(content)
			case prefixError:
				test.error = strings.NewReader(content)
			case prefixCode:
				test.code = bytes.NewReader([]byte(content))
			default:
				return errors.New("Invalid jslt file")
			}
		}

		content = ""
		current = next

		return nil
	}

	closeTest := func() {
		closeSection("")
		if test.isValid() {
			tests = append(tests, test)
		}

		test = testCase{}
	}

	for true {
		if bytesRead, err := fileReader.ReadString('\n'); err == io.EOF {
			closeTest()
			break
		} else if err != nil {
			return nil, err
		} else {
			line := string(bytesRead)
			next := ""

			if strings.HasPrefix(line, prefixTest) {
				closeTest()
				test.name = line[2 : len(bytesRead)-1]
				continue
			} else if strings.HasPrefix(line, prefixInput) {
				next = prefixInput
			} else if strings.HasPrefix(line, prefixOutput) {
				next = prefixOutput
			} else if strings.HasPrefix(line, prefixError) {
				next = prefixError
			} else if strings.HasPrefix(line, prefixCode) {
				next = prefixCode
			}

			if len(next) > 0 {
				if err := closeSection(next); err != nil {
					return nil, err
				}
			} else {
				content = content + line
			}
		}
	}

	return tests, nil
}
