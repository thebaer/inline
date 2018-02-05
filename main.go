package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

type (
	Config struct {
		OutputFile  string
		PackageName string

		Files []string
	}
	inlineFile struct {
		Name     string
		Contents []byte
	}
)

func (f *inlineFile) ContentsString() string {
	return string(f.Contents)
}

func main() {
	cfg := &Config{}
	flag.StringVar(&cfg.OutputFile, "o", "", "Output filename; sent to stdout if none given")
	flag.StringVar(&cfg.PackageName, "p", "main", "Package name")
	flag.Parse()

	cfg.Files = flag.Args()

	if err := Run(cfg); err != nil {
		log.Fatal(err)
	}
}

func Run(cfg *Config) error {
	w := new(bytes.Buffer)
	genText, err := create(cfg)
	if nil != err {
		return fmt.Errorf("failed to expand autogenerated code: %s", err)
	}
	if _, err := w.Write(genText); err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}

	out := os.Stdout
	if cfg.OutputFile != "" {
		if out, err = os.Create(cfg.OutputFile); err != nil {
			return err
		}
	}

	if _, err := w.WriteTo(out); err != nil {
		return err
	}

	if cfg.OutputFile != "" {
		return out.Close()
	}
	return nil
}

func create(cfg *Config) ([]byte, error) {
	tmplParams := struct {
		Invocation, PackageName string

		Files []inlineFile
	}{
		Invocation:  strings.Join(os.Args[1:], " "),
		PackageName: cfg.PackageName,
	}

	files := []inlineFile{}
	// Read in files
	for _, f := range cfg.Files {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v. Skipping.\n", err)
			continue
		}

		files = append(files, inlineFile{
			Name:     f,
			Contents: data,
		})
	}
	tmplParams.Files = files

	t, err := template.New("").Parse(fileTemplate)
	if nil != err {
		return nil, err
	}

	var b bytes.Buffer
	err = t.Execute(&b, tmplParams)
	if nil != err {
		return nil, err
	}

	return b.Bytes(), nil
}

const fileTemplate = `// Code generated by "inline {{.Invocation}}" -- DO NOT EDIT --

package {{.PackageName}}

import (
	"fmt"
	"io/ioutil"
)

func ReadAsset(file string, useLocal bool) ([]byte, error) {
	if useLocal {
		return ioutil.ReadFile(file)
	}
	if f, ok := files[file]; ok {
		return []byte(f), nil
	}
	return nil, fmt.Errorf("file doesn't exist.")
}

var files = map[string]string{
	{{range .Files}}"{{.Name}}": ` + "`{{.ContentsString}}`" + `,{{end}}
}`