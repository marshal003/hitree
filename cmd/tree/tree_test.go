package tree_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/marshal003/hitree/cmd/tree"
)

// Example of running hitree command on a directory having structure as
// a
// 	b
// 		.hidden
// 		normalfile
// 	c
// 		d
// 			e
// 			normal
// 		normal
// normal
func ExampleTree_Print() {
	cleaner, _, root := SetupTestDir("Root")
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

// Only print the directory names
func ExampleTree_Print_dirOnly() {
	cleaner, _, root := SetupTestDir("RootA")
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

// Include hidden files and directory in the result
func ExampleTree_Print_includeHidden() {
	cleaner, _, root := SetupTestDir("RootB")
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

// Don't print report(count of directory and count of file) in the result
func ExampleTree_Print_noReport() {
	cleaner, _, root := SetupTestDir("RootC")
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

// Prune empty directory from result
func ExampleTree_Print_prune() {
	cleaner, _, root := SetupTestDir("RootD")
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

// Inlude files only till certain depth
func ExampleTree_Print_maxLevel() {
	cleaner, _, root := SetupTestDir("RootE")
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
func ExampleTree_Print_includePattern() {
	cleaner, _, root := SetupTestDir("RootF")
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

// Exclude files matching pattern
func ExampleTree_Print_excludePatternWithDepth() {
	cleaner, _, root := SetupTestDir("RootG")
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

type Cleaner func()

func SetupTestDir(root string) (Cleaner, tree.Options, string) {
	tree.InitAurora(false)
	opt := tree.DefaultOptions()
	pwd, _ := os.Getwd()
	tempDir := os.TempDir()
	tempDir = filepath.Join(tempDir, root)
	os.RemoveAll(tempDir)
	os.MkdirAll(filepath.Join(tempDir, "a", "b"), 0777)
	os.MkdirAll(filepath.Join(tempDir, "a", "c", "d", "e"), 0777)
	os.OpenFile(filepath.Join(tempDir, "a", "b", ".hidden"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "a", "b", "normal.go"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "a", "c", "d", "normal.py"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "a", "c", "normal.go"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "a", "normal.py"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "normal.go"), os.O_RDONLY|os.O_CREATE, 0666)
	os.Chdir(tempDir)
	return func() {
		os.Chdir(pwd)
		os.RemoveAll(tempDir)
	}, opt, tempDir
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
