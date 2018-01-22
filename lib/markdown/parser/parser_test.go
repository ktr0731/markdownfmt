package parser

import (
	"testing"

	"github.com/ktr0731/markdownfmt/lib/markdown/token"
	"github.com/stretchr/testify/assert"
)

func Test_readChar(t *testing.T) {
	cases := []struct {
		in string
	}{
		{"makise kurisu"},
		{"ヽ(*ﾟдﾟ)ノ"},
	}

	for _, c := range cases {
		l := New(c.in)
		var err error
		for err != EOS {
			assert.Equal(t, l.c, l.in[l.pos-1])
			err = l.readChar()
		}
	}
}

func Test_readWhileEOL(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{"this is a line", "this is a line"},
		{"line and \n", "line and "},
	}

	for _, c := range cases {
		l := New(c.in)
		actual := l.readWhileEOL()
		assert.Equal(t, c.expected, actual)
	}
}

func Test_readWhileSpace(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{"this is a line", "this"},
		{"line\nand", "line\nand"},
	}

	for _, c := range cases {
		l := New(c.in)
		actual := l.readWhileSpace()
		assert.Equal(t, c.expected, actual)
	}
}

func Test_skipSpace(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{"   foo", "foo"},
		{" bar ", "bar "},
	}

	for _, c := range cases {
		l := New(c.in)
		l.skipSpace()
		actual := l.readWhileEOL()
		assert.Equal(t, c.expected, actual)
	}
}

func Test_readHeading(t *testing.T) {
	ph := func(t token.Heading) *token.Heading { return &t }
	pp := func(t token.Paragraph) *token.Paragraph { return &t }

	cases := map[string]struct {
		in       string
		expected token.Token
	}{
		"normal":        {"### header3", ph(token.Heading{Level: 3, Content: "header3"})},
		"rather than 6": {"####### paragraph", pp(token.Paragraph{Content: "####### paragraph"})},
		"normal:closer": {"# STEINS;GATE #", ph(token.Heading{Level: 1, Content: "STEINS;GATE"})},
	}

	for _, c := range cases {
		l := New(c.in)
		actual := l.readHeading()
		assert.Exactly(t, c.expected, actual)
	}
}
