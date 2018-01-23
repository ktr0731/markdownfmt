package markdown

import (
	"fmt"
	"io"
	"strings"

	bf "gopkg.in/russross/blackfriday.v2"
)

type Printer struct {
	out  io.Writer
	node *bf.Node
}

func NewPrinter(out io.Writer, node *bf.Node) *Printer {
	return &Printer{
		out:  out,
		node: node,
	}
}

func (p *Printer) Fprint() {
	p.printNode()
}

func (p *Printer) print(s string) {
	fmt.Fprint(p.out, s)
}

func (p *Printer) printNode() {
	p.node.Walk(p.visitor)
}
func (p *Printer) heading(h bf.HeadingData) {
	p.print(strings.Repeat("#", h.Level))
}

func (p *Printer) visitor(n *bf.Node, entering bool) bf.WalkStatus {
	if entering {
		p.entering(n)
	} else {
		p.exiting(n)
	}
	return bf.GoToNext
}

func (p *Printer) entering(n *bf.Node) {
	switch n.Type {
	case bf.Heading:
		p.heading(n.HeadingData)
		p.print(" ")
	case bf.HorizontalRule:
		p.print(formatHorizontalRule())
	case bf.Text:
		p.print(formatText(n.Literal))
	case bf.Document:
	case bf.Paragraph:
	default:
		p.print(n.String())
	}
}

// only blocks
// paragraphs, block quotations, lists, headings, rules, and code blocks
func (p *Printer) exiting(n *bf.Node) {
	switch n.Type {
	case bf.Heading:
		p.print("\n")
	case bf.Paragraph:
		p.print("  \n")
	}
}

// rules:
// - cannot use spaces rather than one ("  this   is a text  " -> "this is a text")
// - if p[i] is "\n", replace it by " "
func formatText(p []byte) string {
	var i, bufPos int
	var hasSpce bool
	buf := make([]byte, len(p))
	for ; ; i++ {
		if i >= len(p) {
			return strings.TrimSpace(string(buf[:bufPos]))
		}
		if p[i] == ' ' && hasSpce {
			continue
		} else if p[i] == ' ' {
			hasSpce = true
		} else if p[i] != ' ' && hasSpce {
			hasSpce = false
		}

		if p[i] == '\n' {
			p[i] = ' '
		}

		buf[bufPos] = p[i]
		bufPos++
	}
}

// rules:
// - insert two breaks between ---
func formatHorizontalRule() string {
	return "\n---\n\n"
}
