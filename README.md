# inline

[![GoDoc](https://godoc.org/github.com/thebaer/inline?status.svg)](https://godoc.org/github.com/thebaer/inline)

inline embeds files into Go programs and provides an interface for reading them. It is meant as a drop-in replacement for `ioutil.ReadFile()`.

Usage:

`inline [-o filename.go] [-p packagename] filenames...`

Options:

```
-o
	Output filename. Results go to stdout if empty.
-p="main"
	Package name for the generated file.
```

## Go Generate

inline can be invoked by `go generate`:

`//go:generate inline -o files.go -p main list.txt names.txt`

## Example

First generate code with `inline -o files.go -p main list.txt names.txt`

Then use `ReadAsset(string, bool)` to get the file you want.

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := ReadAsset("list.txt", false)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to read file: %v\n", err)
        os.Exit(1)
    }
	fmt.Printf("%s", data)
}
```
