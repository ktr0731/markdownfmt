package markdown

import blackfriday "gopkg.in/russross/blackfriday.v2"

func Parse(in []byte) *blackfriday.Node {
	return blackfriday.New().Parse(in)
}
