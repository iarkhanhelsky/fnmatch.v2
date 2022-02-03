package fnmatch_test

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/iarkhanhelsky/fnmatch.v2"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

type testCase struct {
	Pattern string   `yaml:"pattern"`
	Input   string   `yaml:"input"`
	Flags   []string `yaml:"flags"`
	Want    bool     `yaml:"want"`
}

func (tc testCase) assert(t *testing.T) {
	r := fnmatch.Match(tc.Pattern, tc.Input, tc.flagMap())
	assert.Equal(t, tc.Want, r)
}

func (tc testCase) flagMap() int {
	if tc.Flags == nil {
		return 0
	}
	flags := 0
	for _, f := range tc.Flags {
		switch f {
		case "fnmatch.FNM_NOESCAPE":
			flags |= fnmatch.FNM_NOESCAPE
		case "fnmatch.FNM_PATHNAME":
			flags |= fnmatch.FNM_PATHNAME
		case "fnmatch.FNM_PERIOD":
			flags |= fnmatch.FNM_PERIOD
		case "fnmatch.FNM_LEADING_DIR":
			flags |= fnmatch.FNM_LEADING_DIR
		case "fnmatch.FNM_CASEFOLD":
			flags |= fnmatch.FNM_CASEFOLD
		case "fnmatch.FNM_IGNORECASE":
			flags |= fnmatch.FNM_IGNORECASE
		case "fnmatch.FNM_FILE_NAME":
			flags |= fnmatch.FNM_FILE_NAME
		}
	}
	return flags
}

func setupTest(t *testing.T) map[string][]testCase {
	files, err := ioutil.ReadDir("testdata/bsd")
	if err != nil {
		t.Error(err)
	}

	testcases := make(map[string][]testCase)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") {
			cases, err := readTestCases(t, path.Join("testdata/bsd", file.Name()))
			if err != nil {
				t.Error(err)
			}

			testcases[file.Name()] = cases
		}
	}
	return testcases
}

func readTestCases(t *testing.T, path string) ([]testCase, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	var cases []testCase
	err = yaml.Unmarshal(bytes, &cases)
	if err != nil {
		t.Error(err)
	}

	return cases, nil
}

func TestAll(t *testing.T) {
	testCases := setupTest(t)
	for name, cases := range testCases {
		for i, tc := range cases {
			t.Run(fmt.Sprintf("%s-%d", name, i), func(t *testing.T) {
				tc.assert(t)
			})
		}
	}
}

func TestManual(t *testing.T) {
	testCase{Pattern: "a/*", Input: "a/bc", Want: true}.assert(t)
	testCase{Pattern: "*", Input: "", Want: true}.assert(t)
	testCase{Pattern: "?", Input: "", Want: false}.assert(t)
}
