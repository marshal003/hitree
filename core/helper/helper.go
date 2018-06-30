package helper

import (
	"os"
	"path/filepath"

	tree "github.com/marshal003/hitree/core"
)

// Cleaner helper utility being used in testcases
type Cleaner func()

// SetupTestDir Helper utility to setup test directory in temp path and then
// return a cleaner function which should be invoked from test case to cleanup
// the temp directory. Testcases needs to ensure to pass unique
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

// SetupFlattenTestDir Helper utility to setup test directory in temp path and then
// return a cleaner function which should be invoked from test case to cleanup
// the temp directory. Testcases needs to ensure to pass unique
func SetupFlattenTestDir(root string) (Cleaner, tree.Options, string) {
	tree.InitAurora(false)
	opt := tree.DefaultOptions()
	pwd, _ := os.Getwd()
	tempDir := os.TempDir()
	tempDir = filepath.Join(tempDir, root)
	os.RemoveAll(tempDir)
	os.MkdirAll(filepath.Join(tempDir, "a"), 0777)
	os.OpenFile(filepath.Join(tempDir, ".hidden"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "golang.go"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "python.py"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "python.pyc"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "normal.py"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "normal.go"), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(filepath.Join(tempDir, "normal.txt"), os.O_RDONLY|os.O_CREATE, 0666)
	os.Chdir(tempDir)
	return func() {
		os.Chdir(pwd)
		os.RemoveAll(tempDir)
	}, opt, tempDir
}
