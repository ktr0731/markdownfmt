package lexer

import (
	"errors"
	"go/token"
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

func (l *Lexer) NextToken() *token.Token {
	return nil
}
