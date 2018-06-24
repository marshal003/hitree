// Package tree implements the data models and helper utilities for
// printing tree structure of the specified directory
//
// By Vinit Kumar Rai <vinitrai.marshal@gmail.com>
package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Stats Data Model to hold the File & Dir Count in the tree.
// Going further this will be enhanced to keep other stats like
// disk utilization, modification time etc at file and dir level
// which could be helpful in including / excluding files and directories based on
// these values.
type Stats struct {
	DirCount  int
	FileCount int
}

// Tree Data Model representing a file/directory as Root and Its childrens(in case of directory).
// It also contains stats of the Root, which currently only have count of all directories and files
// in the Root. Going further we can ustilize Stat to include few more details like space used etc.
type Tree struct {
	Root      os.FileInfo
	Childrens []Tree
	Stats     Stats
}

// String Implements String method of Stringer interface, helpful in debugging.
// For Printing the tree structure like linux command use Print method
func (tree Tree) String() string {
	name := tree.Root.Name()
	childrens := make([]string, len(tree.Childrens))
	for i, t := range tree.Childrens {
		childrens[i] = t.String()
	}
	return fmt.Sprintf("%s->%+v", name, childrens)
}

// Print is a utility method on the tree to print tree structure on console.
// Later we will few more mthods on tree which will allow to output tree result
// to other means like file or socket etc.
func (tree Tree) Print(opt Options) {
	tree.printTree(opt, 0, false)
	if !opt.NoReport {
		fmt.Printf("\n%d directories, %d files\n", tree.Stats.DirCount, tree.Stats.FileCount)
	}
}

//printTree Helper private method to recursively print tree on console
func (tree Tree) printTree(opt Options, padding int, isLastChild bool) {
	tree.printNode(opt)
	l := len(tree.Childrens)
	for index, subtree := range tree.Childrens {
		i := 0
		if canPrune(subtree, opt) {
			continue
		}
		for i < padding {
			if i+1 == padding && isLastChild {
				fmt.Printf("%s", " ")
			} else {
				fmt.Printf("%s", opt.PipeColor("│"))
			}
			fmt.Printf("%s", strings.Repeat(" ", 2))
			i = i + 1
		}

		if l == index+1 {
			fmt.Printf("%s", opt.LLinkColor("└──"))
		} else {
			fmt.Printf("%s", opt.TLinkColor("├──"))
		}
		subtree.printTree(opt, padding+1, index+1 == l)
	}
}

//getColor Helper method to get color based on file type
func (tree Tree) getColor(opt Options) Colorize {
	if tree.Root.IsDir() {
		return opt.DirColor
	}
	return opt.FileColor
}

//printNode Helper private method to print node(Root of tree)
func (tree Tree) printNode(opt Options) {
	colorize := tree.getColor(opt)
	if !opt.ShowFullPath {
		fmt.Printf("%s\n", colorize(tree.Root.Name()))
		return
	}
	path, err := filepath.Abs(tree.Root.Name())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", colorize(path))
}

//canPrune Helper private method to check if the tree can be pruned from output
//Useful for pruning empty directory from output.
func canPrune(tree Tree, opt Options) bool {
	if opt.Prune && tree.Root.IsDir() && tree.Stats.FileCount == 0 {
		return true
	}
	return false
}
