package main

import (
	"io/ioutil"
	"os"

	blackfriday "gopkg.in/russross/blackfriday.v2"

	"github.com/k0kubun/pp"
	"github.com/ktr0731/markdownfmt/lib/markdown"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	nodes := markdown.Parse(b)
	nodes.Walk(visitor)
}

func visitor(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Heading:
		pp.Println(node.HeadingData)
	}
	return blackfriday.GoToNext
}
