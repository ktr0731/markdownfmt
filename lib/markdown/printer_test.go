package markdown

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func read(t *testing.T, fname string) []byte {
	b, err := ioutil.ReadFile(fname)
	require.NoError(t, err)
	return b
}

func TestPrinter(t *testing.T) {
	node := Parse(read(t, "testdata/test1.md"))
	p := NewPrinter(os.Stderr, node)
	p.Fprint()
}

func Test_formatParagraph(t *testing.T) {
	cases := []struct {
		in       string
		expected string
	}{
		{"foo   ", "foo"},
		{" bar ", "bar"},
		{"this   is   a text     ", "this is a text"},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, formatText([]byte(c.in)))
	}
}
