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
	case blackfriday.Document:
	case blackfriday.BlockQuote:
	case blackfriday.List:
	case blackfriday.Item:
	case blackfriday.Paragraph:
		// pp.Println(node.String())
		// return blackfriday.GoToNext
	case blackfriday.Heading:
		// == or #
		pp.Println(node.HeadingData)
		return blackfriday.GoToNext
	case blackfriday.HorizontalRule:
	case blackfriday.Emph:
	case blackfriday.Strong:
	case blackfriday.Del:
	case blackfriday.Link:
	case blackfriday.Image:
	case blackfriday.Text:
	case blackfriday.HTMLBlock:
	case blackfriday.CodeBlock:
	case blackfriday.Softbreak:
	case blackfriday.Code:
	case blackfriday.HTMLSpan:
	case blackfriday.Table:
	case blackfriday.TableCell:
	case blackfriday.TableHead:
	case blackfriday.TableBody:
	case blackfriday.TableRow:
	}
	pp.Println(node.String())
	return blackfriday.GoToNext
}
