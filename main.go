package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/ktr0731/markdownfmt/lib/markdown"
	"github.com/mattn/go-isatty"
)

var (
	write  = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff = flag.Bool("d", false, "display diffs instead of rewriting files")
)

func main() {
	flag.Parse()

	if isatty.IsTerminal(os.Stdin.Fd()) {
		for _, fname := range flag.Args() {
			b, err := readFile(fname)
			if err != nil {
				panic(err)
			}
			format(b, fname)
		}
	} else {
		b := new(bytes.Buffer)
		b.ReadFrom(os.Stdin)
		p := markdown.NewPrinter(os.Stdout, markdown.Parse(b.Bytes()))
		p.Fprint()
	}
}

// if out is os.Stdout, dest is not used.
func format(in []byte, dest string) error {
	var out io.Writer = os.Stdout
	if *write {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer f.Close()
		out, err = os.Create(dest)
		if err != nil {
			return err
		}
		defer out.(*os.File).Close()
	}

	p := markdown.NewPrinter(out, markdown.Parse(in))
	p.Fprint()

	return nil
}

func readFile(fname string) ([]byte, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := new(bytes.Buffer)
	_, err = io.Copy(b, f)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
