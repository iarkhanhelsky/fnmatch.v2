package fnmatch

import (
	"go.uber.org/goleak"
	"io/ioutil"
	"path"
	"testing"
)

func TestAll(t *testing.T) {
	testCases, err := setupTest()
	if err != nil {
		t.Error(err)
	}
	for _, tc := range testCases {
		t.Run(tc.name(), func(t *testing.T) {
			tc.assert(t)
		})
	}
}

func TestManual(t *testing.T) {
	newTestcase("*?", "a", true).assert(t)
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func setupTest() ([]testCase, error) {
	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		return nil, err
	}

	testcases := make([]testCase, 0)
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		dir := f.Name()
		cases, err := readDir(path.Join("testdata", dir), dir)
		if err != nil {
			return nil, err
		}
		testcases = append(testcases, cases...)
	}

	return testcases, nil
}

func BenchmarkMatch(b *testing.B) {
	testcases, err := setupTest()
	if err != nil {
		b.Error(err)
	}
	for _, t := range testcases {
		if t.Skip {
			continue
		}

		var result = false
		b.Run(t.name(), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				result = Match(t.Pattern, t.Input, t.flagMap())
			}
		})
		if result {

		}
	}
}
