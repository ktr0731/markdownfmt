package markdown

import (
	"bytes"
	"io/ioutil"
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
	buf := new(bytes.Buffer)
	p := NewPrinter(buf, node)
	p.Fprint()

	expected := `# header1

---

## header2
this is a text.  

two line text and paragraph.  

many breaks are trimed by AST.  

### header
`

	assert.Equal(t, expected, buf.String())
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
