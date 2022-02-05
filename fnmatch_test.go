package fnmatch

import (
	"io/ioutil"
	"path"
	"testing"
)

func TestAll(t *testing.T) {
	testCases := setupTest(t)
	for _, tc := range testCases {
		t.Run(tc.name(), func(t *testing.T) {
			tc.assert(t)
		})
	}
}

func TestManual(t *testing.T) {
	newTestcase("**/d", "a/b/c/d", true, FNM_PATHNAME).assert(t)
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
		cases, err := readDir(path.Join("testdata", dir), dir)
		if err != nil {
			t.Error(err)
		}
		testcases = append(testcases, cases...)
	}

	return testcases
}
