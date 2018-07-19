// Package core implements the data models and helper utilities for
// printing tree structure of the specified directory
//
// By Vinit Kumar Rai <vinitrai.marshal@gmail.com>
package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Stats Data Model to hold the File & Dir Count in the tree.
// Going further this will be enhanced to keep other stats like
// disk utilization, modification time etc at file and dir level
// which could be helpful in including / excluding files and directories based on
// these values.
type Stats struct {
	DirCount         int       `json:"dir_count"`
	FileCount        int       `json:"file_count"`
	Size             int64     `json:"size"`
	ModificationTime time.Time `json:"mod_time"`
	Permission       string    `json:"permission"`
}

//NewEmptyStats ...
func NewEmptyStats(fi os.FileInfo) Stats {
	stats := Stats{
		DirCount:         0,
		FileCount:        0,
		Size:             fi.Size(),
		ModificationTime: fi.ModTime(),
		Permission:       fi.Mode().Perm().String(),
	}
	return stats
}

// Tree Data Model representing a file/directory as Root and Its childrens(in case of directory).
// It also contains stats of the Root, which currently only have count of all directories and files
// in the Root. Going further we can ustilize Stat to include few more details like space used etc.
type Tree struct {
	Root      os.FileInfo
	Childrens []Tree
	Stats     Stats
}

//JSONTree Json Representation of Tree
type JSONTree struct {
	Name    string     `json:"name"`
	FType   FileType   `json:"file_type"`
	FStats  Stats      `json:"stats"`
	SubTree []JSONTree `json:"subtree"`
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
func (tree Tree) Print(w io.Writer, opt Options) {
	tree.printTree(w, opt, 0, false)
	if !opt.NoReport {
		fmt.Fprintf(w, "\n%d directories, %d files\n", tree.Stats.DirCount, tree.Stats.FileCount)
	}
}

//printTree Helper private method to recursively print tree on console
func (tree Tree) printTree(w io.Writer, opt Options, padding int, isLastChild bool) {
	tree.printNode(w, opt)
	l := len(tree.Childrens)
	for index, subtree := range tree.Childrens {
		i := 0
		if canPrune(subtree, opt) {
			continue
		}
		for i < padding {
			if i+1 == padding && isLastChild {
				fmt.Fprintf(w, "%s", " ")
			} else {
				fmt.Fprintf(w, "%s", opt.PipeColor("│"))
			}
			fmt.Fprintf(w, "%s", strings.Repeat(" ", 2))
			i = i + 1
		}

		if l == index+1 {
			fmt.Fprintf(w, "%s", opt.LLinkColor("└──"))
		} else {
			fmt.Fprintf(w, "%s", opt.TLinkColor("├──"))
		}
		subtree.printTree(w, opt, padding+1, index+1 == l)
	}
}

//getColor Helper method to get color based on file type
func (tree Tree) getColor(opt Options) Colorize {
	if tree.Root.IsDir() {
		return opt.DirColor
	}
	return opt.FileColor
}

//GetExtra ...
func GetExtra(tree Tree, opt Options) string {
	extra := make([]string, 0)
	if opt.PrintUID {
		extra = append(extra, fmt.Sprintf("%d", tree.Root.Sys().(*syscall.Stat_t).Uid))
	}
	if opt.PrintGID {
		extra = append(extra, fmt.Sprintf("%d", tree.Root.Sys().(*syscall.Stat_t).Gid))
	}
	if opt.PrintSize {
		extra = append(extra, fmt.Sprintf("%d", tree.Stats.Size))
	}
	if opt.PrintModTime {
		extra = append(extra, tree.Stats.ModificationTime.Format(opt.TimeFormat))
	}

	if opt.PrintProtection {
		extra = append(extra, tree.Stats.Permission)
	}
	res := strings.Join(extra, " ")
	if len(extra) >= 1 {
		return fmt.Sprintf("[ %s ]", res)
	}
	return ""
}

//printNode Helper private method to print node(Root of tree)
func (tree Tree) printNode(w io.Writer, opt Options) {
	colorize := tree.getColor(opt)
	path, err := tree.NodeName(opt)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%s%s\n", colorize(GetExtra(tree, opt)), colorize(path))
}

//NodeName Get NodeName of the tree
func (tree Tree) NodeName(opt Options) (string, error) {
	if !opt.ShowFullPath {
		return tree.Root.Name(), nil
	}
	return filepath.Abs(tree.Root.Name())
}

//canPrune Helper private method to check if the tree can be pruned from output
//Useful for pruning empty directory from output.
func canPrune(tree Tree, opt Options) bool {
	if opt.Prune && tree.Root.IsDir() && tree.Stats.FileCount == 0 {
		return true
	}
	return false
}

//AsJSONTree Utility function to return tree as json string
func (tree Tree) AsJSONTree(opt Options) JSONTree {
	name, _ := tree.NodeName(opt)
	fileType := GetFileType(tree.Root)
	jsonTree := JSONTree{Name: name, FType: fileType, SubTree: []JSONTree{}}
	if opt.JSONIncludeStats {
		jsonTree.FStats = tree.Stats
	}
	for _, subtree := range tree.Childrens {
		if canPrune(subtree, opt) {
			continue
		}
		child := subtree.AsJSONTree(opt)
		jsonTree.SubTree = append(jsonTree.SubTree, child)
	}
	return jsonTree
}

//AsJSONString ...
func (tree Tree) AsJSONString(opt Options) ([]byte, error) {
	jsonTree := tree.AsJSONTree(opt)
	return json.MarshalIndent(jsonTree, strings.Repeat(" ", int(opt.Indent)), strings.Repeat(" ", int(opt.Indent)))
}
