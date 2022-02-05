# fnmatch.v2

This is fnmatch implementation inspired by [github.com/danwakefield/fnmatch](https://github.com/danwakefield/fnmatch)
and based on Ruby `File.fnmatch` implementation. 

## Installation

```
go get github.com/iarkhanhelsky/fnmatch.v2
```

## Features

| Feature                        | fnmatch.v2         | golang             | fnmatch            | ruby | 
|--------------------------------|--------------------|--------------------|--------------------|------|
| `*`                            | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | 
| `?`                            | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Character range (i.e. `[A-Z]`) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | 
| Windows separator              |                    | :white_check_mark: |                    | :white_check_mark: |
| Recursive match                | :white_check_mark: |                    |                    | :white_check_mark: |
| `FNM_NOESCAPE`                 | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_PATHNAME`                 | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_PERIOD`                   | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_LEADING_DIR`              | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_CASEFOLD`                 | :white_check_mark: |                    | :white_check_mark: | :white_check_mark: |
| `FNM_EXTGLOB`                 |                    |                    |                    | :white_check_mark: |
| `FNM_SHORTNAME` (Windows)      |                    |                    |                    | :white_check_mark: |
| `FNM_SYSCASE`                  |                    |                    |                    | :white_check_mark: |


## Notes 

* https://research.swtch.com/glob