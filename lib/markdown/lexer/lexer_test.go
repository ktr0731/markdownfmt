package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Run("readChar", func(t *testing.T) {
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
	})
}
