# fnmatch.v2
[![Build Status](https://github.com/iarkhanhelsky/fnmatch.v2/actions/workflows/go.yaml/badge.svg)](https://github.com/iarkhanhelsky/fnmatch.v2/actions?workflow=go)
[![Coverage Status](https://coveralls.io/repos/github/iarkhanhelsky/fnmatch.v2/badge.svg?branch=main)](https://coveralls.io/github/iarkhanhelsky/fnmatch.v2?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/iarkhanhelsky/fnmatch.v2.svg)](https://pkg.go.dev/github.com/iarkhanhelsky/fnmatch.v2)

This is fnmatch implementation inspired by [github.com/danwakefield/fnmatch](https://github.com/danwakefield/fnmatch)
and heavily based on Ruby `File.fnmatch`. 

## Installation

```
go get github.com/iarkhanhelsky/fnmatch.v2
```

## Example

```go
import "github.com/iarkhanhelsky/fnmatch.v2"

func main() {
	if fnmatch.Match("foo/*", "foo/bar") {
		println("Matched!")
	}
}
```

## Features

| Feature                        | fnmatch.v2         | golang             | fnmatch            | ruby               | 
|--------------------------------|--------------------|--------------------|--------------------|--------------------|
| `*`                            | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | 
| `?`                            | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Character range (i.e. `[A-Z]`) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | 
| Windows separator              |                    | :white_check_mark: |                    | :white_check_mark: |
| Recursive match                | :white_check_mark: |                    |                    | :white_check_mark: |
| `FNM_NOESCAPE`                 | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_PATHNAME`                 | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_PERIOD`                   | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_LEADING_DIR`              |                    |                    | :white_check_mark: |                    |
| `FNM_CASEFOLD`                 | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_EXTGLOB`                  | :white_check_mark: |                    |                    | :white_check_mark: |
| `FNM_SHORTNAME` (Windows)      |                    |                    |                    | :white_check_mark: |
| `FNM_SYSCASE`                  |                    |                    |                    | :white_check_mark: |


## Notes 

* https://research.swtch.com/glob
* https://www.gnu.org/software/libc/manual/html_node/Wildcard-Matching.html
