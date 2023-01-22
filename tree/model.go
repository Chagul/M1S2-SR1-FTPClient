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

// InitNode Allocate the space for the Children array
func (node *Node) InitNode() {
	node.Children = make([]*Node, 0)
}

// AddChild Add a child to the current node
func (node *Node) AddChild(child *Node) {
	node.Children = append(node.Children, child)
}

// AddChildren Add all the child to the current node
func (node *Node) AddChildren(children []*Node) {
	for _, child := range children {
		alreadyPresent := false
		for _, childInNode := range node.Children {
			if childInNode == child {
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			node.Children = append(node.Children, child)
		}
	}
}

// DisplayTree Display the tree of the current node
// If fullPath is true, display complete path of files, else just the name
// If directoryOnly is true, display only the directories, else display both directories and files
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
