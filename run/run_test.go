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

	"fmt"

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

func (t *testCase) isEmpty() bool {
	return len(t.name) == 0 && t.code == nil && t.output == nil && t.error == nil
}

func (t *testCase) isValid() bool {
	return len(t.name) > 0 && t.code != nil && (t.output != nil || t.error != nil)
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

	actual := bytes.NewBufferString("")
	actualError := bytes.NewBufferString("")

	encounteredError := run.Interpret(test.code, test.input, actual, actualError)

	fail := func(msg string) {
		t.Errorf("\"%s\" failed!\n%s", test.name, msg)
	}

	if test.output != nil {
		expected, _ := ioutil.ReadAll(test.output)

		if actual.String() != string(expected) {
			test.code.Seek(0, io.SeekStart)
			code, _ := ioutil.ReadAll(test.code)
			fail(fmt.Sprintf(
				"Expected output:\n%s\nActual output:\n%s\nErrors:\n%s\nInput:\n%s\n",
				expected,
				actual.String(),
				actualError.String(),
				code,
			))
		}
	} else if len(actual.String()) > 0 {
		fail(fmt.Sprintf(
			"Unexpected output: %s",
			actual.String(),
		))
	}

	if test.error == nil && encounteredError {
		fail(fmt.Sprintf(
			"Got error output: %s\nExpected none",
			actualError.String(),
		))
	} else if test.error != nil {
		expectedError, _ := ioutil.ReadAll(test.error)

		if !encounteredError {
			fail(fmt.Sprintf(
				"Expected error output: %s\nBut interpreter returned success",
				expectedError,
			))
		} else if string(expectedError) != actualError.String() {
			fail(fmt.Sprintf(
				"Expected error output: %s\nBut got:\n%s",
				expectedError,
				actualError.String(),
			))
		}
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

	closeTest := func() error {
		closeSection("")
		if !test.isEmpty() {
			if test.isValid() {
				tests = append(tests, test)
			} else {
				return errors.New("Invalid test file")
			}
		}

		test = testCase{}

		return nil
	}

	for true {
		if bytesRead, err := fileReader.ReadString('\n'); err == io.EOF {
			if err := closeTest(); err != nil {
				return nil, err
			}
			break
		} else if err != nil {
			return nil, err
		} else {
			line := string(bytesRead)
			next := ""

			if strings.HasPrefix(line, prefixTest) {
				if err := closeTest(); err != nil {
					return nil, err
				}
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
