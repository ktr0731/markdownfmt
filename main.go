package main

import (
	"io/ioutil"
	"os"

	"github.com/ktr0731/markdownfmt/lib/markdown"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	p := markdown.NewPrinter(os.Stdout, markdown.Parse(b))
	p.Fprint()
}
