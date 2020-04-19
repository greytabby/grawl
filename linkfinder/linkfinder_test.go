package linkfinder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
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
		</div>
	</body>
</html>
`
)

func TestFindLink(t *testing.T) {
	node, err := html.Parse(strings.NewReader(testHTML))
	assert.NoError(t, err)
	links := FindLinks(node)

	assert.Equal(t, 6, len(links))
	assert.NotEmpty(t, links)

	t.Log(links)
}
