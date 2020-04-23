package scrape

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Matcher should return true if you desire to select node.
type Matcher func(*html.Node) bool

// FindAll return all nodes that match matcher function.
func FindAll(node *html.Node, matcher func(*html.Node) bool) []*html.Node {
	nodes := make([]*html.Node, 0)
	if matcher(node) {
		nodes = append(nodes, node)
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		mathed := FindAll(n, matcher)
		nodes = append(nodes, mathed...)
	}
	return nodes
}

// Attr return value of attribute.
func Attr(node *html.Node, key string) string {
	for _, v := range node.Attr {
		if v.Key == key {
			return v.Val
		}
	}
	return ""
}

// Text return joined textdata from all child `html.TextNode` node.
func Text(node *html.Node) string {
	textNodes := FindAll(node, func(n *html.Node) bool {
		return n.Type == html.TextNode
	})

	texts := make([]string, 0)
	for _, tn := range textNodes {
		text := strings.TrimSpace(tn.Data)
		if text == "" {
			continue
		}
		texts = append(texts, text)
	}
	return strings.Join(texts, " ")
}

// ByTag return `Matcher` that match node that have specified tag.
func ByTag(tag atom.Atom) Matcher {
	return func(node *html.Node) bool {
		return node.DataAtom == tag
	}
}
