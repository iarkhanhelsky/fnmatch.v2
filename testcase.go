package fnmatch

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

type testCase struct {
	Pattern string   `yaml:"pattern"`
	Input   string   `yaml:"input"`
	Flags   []string `yaml:"flags"`
	Want    bool     `yaml:"want"`
	Skip    bool     `yaml:"skip"`

	group string
	file  string
	index int
}

func newTestcase(pattern string, input string, want bool, flags ...int) testCase {
	stringflags := []string{}
	for _, f := range flags {
		switch f {
		case FNM_NOESCAPE:
			stringflags = append(stringflags, "fnmatch.FNM_NOESCAPE")
		case FNM_PATHNAME:
			stringflags = append(stringflags, "fnmatch.FNM_PATHNAME")
		case FNM_PERIOD:
			stringflags = append(stringflags, "fnmatch.FNM_PERIOD")
		case FNM_LEADING_DIR:
			stringflags = append(stringflags, "fnmatch.FNM_LEADING_DIR")
		case FNM_CASEFOLD:
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
	r := Match(tc.Pattern, tc.Input, tc.flagMap())
	if r != tc.Want && tc.Skip {
		t.Skip(tc.string())
		return
	}
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
			flags |= FNM_NOESCAPE
		case "fnmatch.FNM_PATHNAME":
			flags |= FNM_PATHNAME
		case "fnmatch.FNM_PERIOD":
			flags |= FNM_PERIOD
		case "fnmatch.FNM_DOTMATCH":
			flags |= FNM_DOTMATCH
		case "fnmatch.FNM_LEADING_DIR":
			flags |= FNM_LEADING_DIR
		case "fnmatch.FNM_CASEFOLD":
			flags |= FNM_CASEFOLD
		case "fnmatch.FNM_IGNORECASE":
			flags |= FNM_IGNORECASE
		case "fnmatch.FNM_FILE_NAME":
			flags |= FNM_FILE_NAME
		default:
			panic(f)
		}
	}
	return flags
}

func readDir(dirPath string, group string) ([]testCase, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var testcases []testCase
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") {
			cases, err := readTestCases(path.Join(dirPath, file.Name()))
			if err != nil {
				return nil, err
			}

			for i, tc := range cases {
				tc.group = group
				tc.file = file.Name()
				tc.index = i
				testcases = append(testcases, tc)
			}
		}
	}

	return testcases, nil
}

func readTestCases(path string) ([]testCase, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cases []testCase
	err = yaml.Unmarshal(bytes, &cases)
	if err != nil {
		return nil, err
	}

	return cases, nil
}