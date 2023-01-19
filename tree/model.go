package tree

import (
	"fmt"
	"strings"
)

type Node struct {
	Filepath    string
	Filename    string
	Children    []*Node
	IsDirectory bool
	Depth       int
}

func (node *Node) InitNode() {
	node.Children = make([]*Node, 0)
}

func (node *Node) AddChild(child *Node) {
	node.Children = append(node.Children, child)
}

func (node *Node) AddChildren(children []*Node) {
	for _, child := range children {
		node.Children = append(node.Children, child)
	}
}

func (node *Node) DisplayTree() {
	space := "    "
	trail := "---"
	//branch := "│   "
	tee := "├── "
	last := "└── "
	if node.IsDirectory {
		if len(node.Children) == 0 {
			fmt.Printf("%s%s\n", last, node.Filename)
		} else {
			fmt.Printf("%s%s\n", tee, node.Filename)
		}
		for _, child := range node.Children {
			fmt.Printf("%s", strings.Repeat(space, child.Depth))
			child.DisplayTree()
		}
	} else {
		fmt.Printf("%s%s\n", trail, node.Filename)
	}
}
