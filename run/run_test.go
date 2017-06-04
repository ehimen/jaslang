package run_test

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Relative to GOPATH
var jsltPath string = "src/github.com/ehimen/jaslang/jslt"

type testCase struct {
	name     string
	code     io.RuneReader
	output   io.Reader
	input    io.Reader
	expected io.Reader
}

func (t *testCase) isValid() bool {
	return len(t.name) > 0
}

func TestJslt(t *testing.T) {
	goPath, exists := os.LookupEnv("GOPATH")

	if !exists {
		t.Log("Skipping JSLT tests as GOPATH is not defined")
		return
	} else {
		t.Log("Running JSLT tests using " + jsltPath)
	}

	tests := []testCase{}

	for _, path := range filepath.SplitList(goPath) {
		path = filepath.Join(path, jsltPath)

		if files, err := filepath.Glob(filepath.Join(path, "*.jslt")); err != nil {
			return
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

	t.Log(tests)
}

func parseFile(path string) ([]testCase, error) {
	tests := []testCase{}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	fileReader := bufio.NewReader(file)

	prefixTest := "<!"
	//prefixInput := "<<<INPUT"
	//prefixOutput := "<<<OUTPUT"
	//prefixCode := "<<<CODE"
	//prefixError := "<<<ERROR"

	test := testCase{}

	for true {
		if bytes, _, err := fileReader.ReadLine(); err == io.EOF {
			if test.isValid() {
				tests = append(tests, test)
			}
			break
		} else if err != nil {
			return nil, err
		} else {
			line := string(bytes)

			if strings.HasPrefix(line, prefixTest) {
				if test.isValid() {
					tests = append(tests, test)
				}
				test = testCase{}
				test.name = line[1:]
			}
		}
	}

	return tests, nil
}
