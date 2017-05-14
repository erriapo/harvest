// Package trie is a very limited radix tree implementation.
//
// For information about radix trees, see https://en.wikipedia.org/wiki/Radix_tree.
package trie

import "strings"

type edge struct {
	node    *node
	label   string
	counter uint
}

type node struct {
	edges []*edge
}

// Trie represents a radix tree.
type Trie struct {
	root *node
}

// WalkFn is used when walking the tree.
type WalkFn func(prefix string, count uint)

func (n *node) appendEdge(e *edge) {
	n.edges = append(n.edges, e)
}

func (n *node) walk(fn WalkFn, prefixes []string) {
	for _, edge := range n.edges {
		newPrefixes := append(prefixes, edge.label)
		fn(
			strings.Join(newPrefixes, ""),
			edge.counter,
		)

		if edge.node != nil {
			edge.node.walk(fn, newPrefixes)
		}
	}
}

func commonPrefix(s1, s2 string) int {
	max := len(s1)
	if l := len(s2); l < max {
		max = l
	}

	var i int
	for i = 0; i < max; i++ {
		if s1[i] != s2[i] {
			break
		}
	}

	return i
}

func (e *edge) split(length int) {
	oldLabel := e.label
	oldNode := e.node

	e.label = oldLabel[:length]
	e.node = &node{}
	e.node.appendEdge(&edge{
		label:   oldLabel[length:],
		counter: e.counter - 1,
		node:    oldNode,
	})
}

// New returns an empty Trie.
func New() *Trie {
	return &Trie{&node{[]*edge{}}}
}

// Insert adds an element to the tree.
func (t *Trie) Insert(s string) {
	var nextEdge *edge
	traverseNode := t.root
	elementsFound := 0

	for traverseNode != nil {
		nextEdge = nil

		for _, e := range traverseNode.edges {
			prefixLength := commonPrefix(s[elementsFound:], e.label)
			if prefixLength == 0 {
				continue
			}

			elementsFound += prefixLength
			e.counter++
			nextEdge = e

			if prefixLength < len(e.label) {
				e.split(prefixLength)
			}

			break
		}

		if nextEdge == nil {
			break
		}

		traverseNode = nextEdge.node
	}

	if elementsFound < len(s) {
		if traverseNode == nil {
			traverseNode = &node{}
			nextEdge.node = traverseNode
		}

		traverseNode.appendEdge(&edge{
			label:   s[elementsFound:],
			counter: 1,
		})
	}
}

// Walk is used to walk the tree edges.
func (t *Trie) Walk(fn WalkFn) {
	t.root.walk(fn, []string{})
}
