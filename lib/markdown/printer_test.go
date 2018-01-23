package markdown

import (
	"bytes"
	"io/ioutil"
	"testing"

	bf "gopkg.in/russross/blackfriday.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func read(t *testing.T, fname string) []byte {
	b, err := ioutil.ReadFile(fname)
	require.NoError(t, err)
	return b
}

func TestPrinter(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
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
	})

	t.Run("test2", func(t *testing.T) {
		node := Parse(read(t, "testdata/test2.md"))
		buf := new(bytes.Buffer)
		p := NewPrinter(buf, node)
		p.Fprint()

		expected := `> 引用
引用2`
		assert.Equal(t, expected, buf.String())
	})
}

func Test_formatText(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		cases := []struct {
			in       string
			expected string
		}{
			{"foo   ", "foo"},
			{" bar ", "bar"},
			{"this   is   a text     ", "this is a text"},
		}

		p := NewPrinter(nil, nil)
		for _, c := range cases {
			assert.Equal(t, c.expected, p.formatText(&bf.Node{Literal: []byte(c.in)}))
		}
	})

	t.Run("in container", func(t *testing.T) {
		cases := []struct {
			in       string
			expected string
		}{
			{"first\n    second", "first\nsecond"},
		}

		p := newPrinterWithinContainer(nil, nil)
		for _, c := range cases {
			assert.Equal(t, c.expected, p.formatText(&bf.Node{Literal: []byte(c.in)}))
		}
	})
}

func Test_breakLine(t *testing.T) {
	in := "this is a long text for testing."
	expected := `this is a
long text
for
testing.`
	actual := breakLine(in, 10)
	assert.Equal(t, expected, actual)
}
