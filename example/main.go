package example

import "github.com/iarkhanhelsky/fnmatch.v2"

func main() {
	if fnmatch.Match("foo/*", "foo/bar") {
		println("Matched!")
	}
}
