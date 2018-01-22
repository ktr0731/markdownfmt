package parser

import (
	"errors"
	"strings"

	"github.com/ktr0731/markdownfmt/lib/markdown/token"
)

var (
	EOS = errors.New("end of source")
)

type Lexer struct {
	in  []rune
	pos int
	c   rune
}

func New(in string) *Lexer {
	l := &Lexer{
		in: []rune(in),
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() error {
	if len(l.in) <= l.pos {
		return EOS
	}
	l.c = l.in[l.pos]
	l.pos++
	return nil
}

func (l *Lexer) readWhile(c ...rune) string {
	var s []rune
	for {
		for _, r := range c {
			if l.c == r {
				return string(s)
			}
		}

		s = append(s, l.c)
		err := l.readChar()

		if err == EOS {
			return string(s)
		} else if err != nil {
			panic(err)
		}
	}
}

func (l *Lexer) readWhileEOL() string {
	return l.readWhile('\n')
}

func (l *Lexer) readWhileSpace() string {
	return l.readWhile(' ')
}

func (l *Lexer) skipSpace() {
	for {
		switch l.c {
		case '\t':
		case ' ':
			if err := l.readChar(); err == EOS {
				return
			} else if err != nil {
				panic(err)
			}
			continue
		default:
			return
		}
	}
}

func (l *Lexer) readHeading() token.Token {
	// read header
	level := 1
L:
	for level <= 6 {
		if err := l.readChar(); err == EOS {
			return &token.Heading{Level: level}
		} else if err != nil {
			panic(err)
		}

		switch l.c {
		case '#':
			level++
		case ' ':
			break L
		case '\n':
			return &token.Heading{Level: level}
		}
	}

	// the line is paragrah
	if level > 6 {
		return &token.Paragraph{Content: strings.Repeat("#", level-1) + strings.TrimSpace(l.readWhileEOL())}
	}

	l.skipSpace()

	// we discard closing sequence
	return &token.Heading{
		Level:   level,
		Content: l.readWhile(' ', '\n'),
	}
}

func (l *Lexer) NextToken() token.Token {
	switch l.c {
	case '#':
		return l.readHeading()
	}
	return nil
}
