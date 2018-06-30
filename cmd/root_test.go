package cmd_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/marshal003/hitree/core/helper"
	"github.com/spf13/cobra"
)

type HiTree cobra.Command

// Example of running hitree command on a directory having structure as
// root
// 	a
// 		b
// 			.hidden
// 			normal.go
// 		c
// 			d
// 				e
// 				normal.py
// 			normal.go
//  	normal.py
// 	normal.go
func ExampleHiTree() {
	cleaner, _, root := helper.SetupTestDir("Root")
	defer cleaner()
	// $ hitree root
	execute("hitree", root)
	// Output:
	// Root
	// ├──a
	// │  ├──b
	// │  │  └──normal.go
	// │  ├──c
	// │  │  ├──d
	// │  │  │  ├──e
	// │  │  │  └──normal.py
	// │  │  └──normal.go
	// │  └──normal.py
	// └──normal.go
	//
	// 5 directories, 5 files
}

// // Only print the directory names
func ExampleHiTree_dirOnly() {
	cleaner, _, root := helper.SetupTestDir("RootA")
	defer cleaner()
	// $ hitree root --dironly
	execute("hitree", root, "--dironly")
	// Output:
	// RootA
	// └──a
	//    ├──b
	//    └──c
	// │     └──d
	// │  │     └──e
	//
	// 5 directories, 5 files
}

// // Include hidden files and directory in the result
func ExampleHiTree_includeHidden() {
	cleaner, _, root := helper.SetupTestDir("RootB")
	defer cleaner()
	// $ hitree root --all
	execute("hitree", root, "--all")
	// Output:
	// RootB
	// ├──a
	// │  ├──b
	// │  │  ├──.hidden
	// │  │  └──normal.go
	// │  ├──c
	// │  │  ├──d
	// │  │  │  ├──e
	// │  │  │  └──normal.py
	// │  │  └──normal.go
	// │  └──normal.py
	// └──normal.go
	//
	// 5 directories, 6 files
}

// // Don't print report(count of directory and count of file) in the result
func ExampleHiTree_noReport() {
	cleaner, _, root := helper.SetupTestDir("RootC")
	defer cleaner()
	// $ hitree root --noreport
	execute("hitree", root, "--noreport")
	// Output:
	// RootC
	// ├──a
	// │  ├──b
	// │  │  └──normal.go
	// │  ├──c
	// │  │  ├──d
	// │  │  │  ├──e
	// │  │  │  └──normal.py
	// │  │  └──normal.go
	// │  └──normal.py
	// └──normal.go
}

// // Prune empty directory from result
func ExampleHiTree_prune() {
	cleaner, _, root := helper.SetupTestDir("RootD")
	defer cleaner()
	// $ hitree root --prune
	execute("hitree", root, "--prune")
	// Output:
	// RootD
	// ├──a
	// │  ├──b
	// │  │  └──normal.go
	// │  ├──c
	// │  │  ├──d
	// │  │  │  ├──e
	// │  │  │  └──normal.py
	// │  │  └──normal.go
	// │  └──normal.py
	// └──normal.go
	//
	// 5 directories, 5 files
}

// // Inlude files only till certain depth
func ExampleHiTree_maxLevel() {
	cleaner, _, root := helper.SetupTestDir("RootE")
	defer cleaner()
	// $ hitree root --level
	execute("hitree", root, "--level=2")
	// Output:
	// RootE
	// ├──a
	// │  ├──b
	// │  ├──c
	// │  └──normal.py
	// └──normal.go
	//
	// 3 directories, 2 files
}

// Inlude files which matches pattern
func ExampleHiTree_includePattern() {
	cleaner, _, root := helper.SetupTestDir("RootF")
	defer cleaner()
	// $ hitree root --includepattern=*.go
	execute("hitree", root, "--includepattern=*.go")
	// Output:
	// RootF
	// ├──a
	// │  ├──b
	// │  │  └──normal.go
	// │  └──c
	// │     ├──d
	// │  │  │  └──e
	// │     └──normal.go
	// └──normal.go
	//
	// 5 directories, 3 files
}

// // Exclude files matching pattern
func ExampleHiTree_patternWithDepth() {
	cleaner, _, root := helper.SetupTestDir("RootG")
	defer cleaner()
	// $ hitree root -I=*.go --level=2
	execute("hitree", root, "-I=*.py", "--level=2")
	// Output:
	// RootG
	// ├──a
	// │  ├──b
	// │  └──c
	// └──normal.go
	//
	// 3 directories, 1 files
}

func execute(command, root string, args ...string) {
	args = append([]string{root, "--nocolor"}, args...)
	path := fmt.Sprintf("PATH=%s:%s", os.Getenv("PATH"), os.Getenv("GOPATH"))
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(), path)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
