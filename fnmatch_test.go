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

	group string
	file  string
	index int
}

func newTestcase(pattern string, input string, want bool, flags ...int) testCase {
	stringflags := []string{}
	for _, f := range flags {
		switch f {
		case fnmatch.FNM_NOESCAPE:
			stringflags = append(stringflags, "fnmatch.FNM_NOESCAPE")
		case fnmatch.FNM_PATHNAME:
			stringflags = append(stringflags, "fnmatch.FNM_PATHNAME")
		case fnmatch.FNM_PERIOD:
			stringflags = append(stringflags, "fnmatch.FNM_PERIOD")
		case fnmatch.FNM_LEADING_DIR:
			stringflags = append(stringflags, "fnmatch.FNM_LEADING_DIR")
		case fnmatch.FNM_CASEFOLD:
			stringflags = append(stringflags, "fnmatch.FNM_CASEFOLD")
		default:
			panic(f)
		}
	}
	return testCase{
		Pattern: pattern, Input: input, Want: want, Flags: stringflags,
	}
}

func (tc testCase) assert(t *testing.T) {
	r := fnmatch.Match(tc.Pattern, tc.Input, tc.flagMap())
	assert.Equal(t, tc.Want, r, tc.string())
}

func (tc testCase) string() string {
	flags := "0"
	if len(tc.Flags) > 0 {
		flags = strings.Join(tc.Flags, " | ")
	}

	return fmt.Sprintf("fnmatch('%s', '%s', %s) -> %t",
		tc.Pattern, tc.Input, flags, tc.Want)
}

func (tc testCase) name() string {
	s := fmt.Sprintf("test_%s_%s_%d", tc.group, tc.file, tc.index)
	s = strings.Replace(strings.Replace(s, "-", "_", -1), ".", "_", -1)
	return s
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
		case "fnmatch.FNM_DOTMATCH":
			flags |= fnmatch.FNM_DOTMATCH
		case "fnmatch.FNM_LEADING_DIR":
			flags |= fnmatch.FNM_LEADING_DIR
		case "fnmatch.FNM_CASEFOLD":
			flags |= fnmatch.FNM_CASEFOLD
		case "fnmatch.FNM_IGNORECASE":
			flags |= fnmatch.FNM_IGNORECASE
		case "fnmatch.FNM_FILE_NAME":
			flags |= fnmatch.FNM_FILE_NAME
		default:
			panic(f)
		}
	}
	return flags
}

func setupTest(t *testing.T) []testCase {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Error(err)
	}

	testcases := make([]testCase, 0)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		dir := f.Name()
		files, err := ioutil.ReadDir(path.Join("testdata", dir))
		if err != nil {
			t.Error(err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".yaml") {
				cases, err := readTestCases(t, path.Join("testdata", dir, file.Name()))
				if err != nil {
					t.Error(err)
				}

				for i, tc := range cases {
					tc.group = dir
					tc.file = file.Name()
					tc.index = i
					testcases = append(testcases, tc)
				}

			}
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
	for _, tc := range testCases {
		t.Run(tc.name(), func(t *testing.T) {
			tc.assert(t)
		})
	}
}

func TestManual(t *testing.T) {
	newTestcase("**/d", "a/b/c/d", true, fnmatch.FNM_PATHNAME).assert(t)
}
