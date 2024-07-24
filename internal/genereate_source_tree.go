package internal

import (
	"path/filepath"
	"strings"
)

type Node struct {
	Name     string
	Children []*Node
}

func GenerateSourceTree(files []FileInfo) *Node {
	root := &Node{Name: "root", Children: []*Node{}}

	for _, file := range files {
		parts := strings.Split(filepath.ToSlash(file.Path), "/")
		current := root

		for _, part := range parts {
			child := findOrCreateChild(current, part)
			current = child
		}
	}

	return root
}

func findOrCreateChild(node *Node, name string) *Node {
	for _, child := range node.Children {
		if child.Name == name {
			return child
		}
	}

	newChild := &Node{Name: name, Children: []*Node{}}
	node.Children = append(node.Children, newChild)
	return newChild
}

func (n *Node) String() string {
	return n.stringify(0)
}

func (n *Node) stringify(depth int) string {
	result := strings.Repeat("  ", depth) + n.Name + "\n"
	for _, child := range n.Children {
		result += child.stringify(depth + 1)
	}
	return result
}