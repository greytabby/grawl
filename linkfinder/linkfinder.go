package linkfinder

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	findAnchor = func(n *html.Node) bool {
		return n.DataAtom == atom.A
	}
)

func FindLinks(node *html.Node) []string {
	nodes := findAll(node, findAnchor)
	if len(nodes) == 0 {
		return []string{}
	}
	
	links := make([]string, 0)
	for _, n := range nodes {
		links = append(links, attr(n, "href"))
	}
	return links
}

func findAll(node *html.Node, matcher func(*html.Node) bool) []*html.Node {
	nodes := make([]*html.Node, 0)
	if matcher(node) {
		nodes = append(nodes, node)
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		mathed := findAll(n, matcher)
		nodes = append(nodes, mathed...)
	}
	return nodes
}

func attr(node *html.Node, key string) string {
	for _, v := range node.Attr {
		if v.Key == key {
			return v.Val
		}
	}
	return ""
}
