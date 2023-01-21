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

func (node *Node) DisplayTree(fullPath bool, directoryOnly bool) {
	space := "    "
	tee := "├── "
	last := "└── "
	var strFile = ""
	if fullPath {
		strFile = node.Filepath
	} else {
		strFile = node.Filename
	}
	if node.IsDirectory {
		repeatSpace := 0
		if node.Depth == 0 {
			fmt.Printf("%s\n", strFile)
		} else {
			repeatSpace = node.Depth + 1
			fmt.Printf("%s\n", strFile)
		}
		firstOnePrinted := false
		for i, child := range node.Children {
			if directoryOnly {
				if child.IsDirectory {
					if i != len(node.Children)-1 || !firstOnePrinted {
						fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), tee)
						firstOnePrinted = true
					} else {
						fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), last)
					}
					child.DisplayTree(fullPath, directoryOnly)
				}
				continue
			}
			if i == len(node.Children)-1 {
				fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), last)
			} else if i == 0 {
				fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), tee)
			} else {
				fmt.Printf("%s%s", strings.Repeat(space, repeatSpace), tee)
			}
			child.DisplayTree(fullPath, directoryOnly)
		}
	} else if !directoryOnly {
		fmt.Printf("%s\n", strFile)
	}
}
