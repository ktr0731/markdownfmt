package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	bf "gopkg.in/russross/blackfriday.v2"
)

type Printer struct {
	out         io.Writer
	node        *bf.Node
	inContainer bool
}

func NewPrinter(out io.Writer, node *bf.Node) *Printer {
	return &Printer{
		out:  out,
		node: node,
	}
}

func newPrinterWithinContainer(out io.Writer, node *bf.Node) *Printer {
	p := NewPrinter(out, node)
	p.inContainer = true
	return p
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
		return p.entering(n)
	}
	return p.exiting(n)
}

func (p *Printer) entering(n *bf.Node) bf.WalkStatus {
	switch n.Type {
	case bf.Heading:
		p.heading(n.HeadingData)
		p.print(" ")
	case bf.HorizontalRule:
		p.print(formatHorizontalRule())
	case bf.Code:
		p.print(formatCode(n))
	case bf.Text:
		p.print(p.formatText(n))
	case bf.BlockQuote:
		p.print(formatBlockQuote(n))
		return bf.SkipChildren
	case bf.Image:
		p.print(formatImage(n))
	case bf.Item:
		p.print(formatItem(n))
	case bf.Strong:
		p.print(formatStrong(n))
		return bf.SkipChildren
	case bf.Hardbreak:
		p.print("  \n")
	case bf.Link:
		p.print(formatLink(n))
		return bf.SkipChildren

	case bf.List:
	case bf.Document:
	case bf.CodeBlock:
	case bf.Paragraph:
	default:
		p.print(n.String())
	}
	return bf.GoToNext
}

// only blocks
// paragraphs, block quotations, lists, headings, rules, and code blocks
func (p *Printer) exiting(n *bf.Node) bf.WalkStatus {
	switch n.Type {
	case bf.Heading:
		p.print("\n")
	case bf.Paragraph:
		p.print(formatParagraphExiting(n))
	case bf.CodeBlock:
	case bf.List:
		if n.Next != nil && n.Next.Type == bf.Paragraph {
			p.print("\n")
		}
	}
	return bf.GoToNext
}

// rules:
// - cannot use spaces rather than one ("  this   is a text  " -> "this is a text")
// - if p[i] is "\n", replace it by " "
//	- however, if the node is within a container, this rule is ignored
// - max width is 80
func (p *Printer) formatText(n *bf.Node) string {
	l := n.Literal
	var i, bufPos int
	var hasSpace, isHead bool
	buf := make([]byte, len(l))
	for ; ; i++ {
		if i >= len(l) {
			return breakLine(strings.TrimSpace(string(buf[:bufPos])), 80)
		}
		if l[i] == ' ' {
			if (p.inContainer && isHead) || hasSpace {
				continue
			}
			hasSpace = true
		} else if l[i] != ' ' && hasSpace {
			hasSpace = false
			isHead = false
		}

		if l[i] == '\n' {
			isHead = true
			if !p.inContainer {
				l[i] = ' '
			}
		}

		buf[bufPos] = l[i]
		bufPos++
	}
}

// rules:
// - insert two breaks between ---
func formatHorizontalRule() string {
	return "\n---\n\n"
}

// rules:
// - insert a break to under the paragraph if next node is a block
func formatParagraphExiting(n *bf.Node) string {
	s := "  \n"
	if n.Next != nil && isBlock(n.Next.Type) {
		s += "\n"
	}
	return s
}

func formatBlockQuote(n *bf.Node) string {
	buf := new(bytes.Buffer)
	if n.FirstChild == nil {
		return "> " // ホンマか
	}
	p := newPrinterWithinContainer(buf, n.FirstChild)
	p.Fprint()
	return "> " + strings.TrimSpace(buf.String())
}

func formatImage(n *bf.Node) string {
	return fmt.Sprintf("![%s](%s)  ", string(n.Literal), string(n.LinkData.Destination))
}

func formatItem(n *bf.Node) string {
	return "- " + string(n.Literal)
}

func formatLink(n *bf.Node) string {
	return fmt.Sprintf("[%s](%s)  ", string(n.FirstChild.Literal), string(n.LinkData.Destination))
}

// rules:
// - insert two spaces between text
func formatStrong(n *bf.Node) string {
	buf := new(bytes.Buffer)
	if n.FirstChild == nil {
		return ""
	}
	p := newPrinterWithinContainer(buf, n.FirstChild)
	p.Fprint()
	return fmt.Sprintf(" **%s** ", buf.String())
}

func formatCode(n *bf.Node) string {
	return fmt.Sprintf("``` %s\n```\n", string(n.Literal))
}

func isBlock(t bf.NodeType) bool {
	return t == bf.Paragraph || t == bf.BlockQuote || t == bf.List || t == bf.Heading || t == bf.CodeBlock
}

func breakLine(s string, w int) string {
	res := [][]string{}
	sp := strings.Split(s, " ")
	l := []string{}
	var cnt int
	for _, p := range sp {
		// last element
		if p == sp[len(sp)-1] {
			if cnt+len(l)+len(p) > w {
				res = append(res, l)
			}
			l = append(l, p)
			res = append(res, l)
			break
		}

		if cnt+len(l)+len(p) > w {
			res = append(res, l)
			l = []string{}
			cnt = 0
		}
		l = append(l, p)
		cnt += len(p)
	}

	tmp := make([]string, len(res))
	for i, a := range res {
		tmp[i] = strings.Join(a, " ")
	}
	return strings.Join(tmp, "\n")
}
