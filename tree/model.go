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
	tee := "├── "
	last := "└── "

	if node.IsDirectory {
		repeatSpace := 0
		if node.Depth == 0 {
			fmt.Printf("%s\n", node.Filename)
		} else {
			repeatSpace = node.Depth + 1
			fmt.Printf("%s\n", node.Filename)
		}
		for i, child := range node.Children {
			if i == len(node.Children)-1 {
				fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), last)
			} else if i == 0 {
				fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), tee)
			} else {
				fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), tee)
			}
			child.DisplayTree()
		}
	} else {
		fmt.Printf("%s\n", node.Filename)
	}
}
