package token

const (
	HEADING = iota
)

type Token interface {
	isToken()
}

type Heading struct {
	Level   int
	Content string
}

func (h *Heading) isToken() {}

type Paragraph struct {
	Content string
}

func (p *Paragraph) isToken() {}
