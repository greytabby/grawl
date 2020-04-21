package linkfinder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	testHTML = `<!DOCTYPE html>
<html>
	<body>
		<h2>HTML Links</h2>
		<div id="a">
			B
			<a href="#jump">Jump</a>
		</div>
		<div id="b">
			<p>test</p>
			A
			<a href="../relative/link">relative</a>
			<div id="bc" class="class1">
				<a href="https://abs.test.link">absolute</a>
			</div>
			<a href="javascript:alert('hello');">javacript</a>
			<a href="ftp://ftp.test.link">javacript</a>
			<a href="">nothing</a>
			<span>
				textA
				<p>textB</p>
				<ul>
					<li>textC</li>
				</ul>
			</span>
			textD
		</div>
	</body>
</html>
`
)

func TestFindAll(t *testing.T) {
	node, err := html.Parse(strings.NewReader(testHTML))
	assert.NoError(t, err)

	got := FindAll(node, ByTag(atom.A))
	assert.Equal(t, 6, len(got))
}

func TestAttr(t *testing.T) {
	node, err := html.Parse(strings.NewReader(testHTML))
	assert.NoError(t, err)

	nodes := FindAll(node, ByTag(atom.Div))
	want := []struct {
		id    string
		class string
	}{
		{"a", ""},
		{"b", ""},
		{"bc", "class1"},
	}

	for i, v := range nodes {
		id := Attr(v, "id")
		class := Attr(v, "class")
		assert.Equal(t, want[i].id, id)
		assert.Equal(t, want[i].class, class)
	}
}

func TestText(t *testing.T) {
	node, err := html.Parse(strings.NewReader(testHTML))
	assert.NoError(t, err)

	spanNode := FindAll(node, ByTag(atom.Span))
	assert.Equal(t, 1, len(spanNode))

	got := Text(spanNode[0])
	want := "textA textB textC"
	assert.Equal(t, want, got)
}
