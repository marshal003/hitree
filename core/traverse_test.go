package core_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/marshal003/hitree/core"
	"github.com/marshal003/hitree/core/helper"
)

// SetupDir returns directory structure
//
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
// normal.go
func TestCountOfRootsChildren(t *testing.T) {
	cleaner, opt, root := helper.SetupTestDir(uuid.New().String())
	defer cleaner()
	tree, error := core.TraverseDir(root, opt, -1)
	if error != nil {
		t.Errorf("Unable to traverse tree rooted at %s", root)
	}
	if len(tree.Childrens) != 2 {
		t.Errorf("Expected to get 2 direct children of but got %d", len(tree.Childrens))
	}
}

func TestDirStat(t *testing.T) {
	cleaner, opt, root := helper.SetupTestDir(uuid.New().String())
	defer cleaner()
	tree, error := core.TraverseDir(root, opt, -1)
	if error != nil {
		t.Errorf("Unable to traverse tree rooted at %s", root)
	}
	if tree.Stats.DirCount != 5 {
		t.Errorf("Expected to get 5 directories but got %d", tree.Stats.DirCount)
	}

	if tree.Stats.FileCount != 5 {
		t.Errorf("Expected to get 5 files but got %d", tree.Stats.FileCount)
	}
}

func TestDirStatIncludingHidden(t *testing.T) {
	cleaner, opt, root := helper.SetupTestDir(uuid.New().String())
	opt.IncludeHidden = true
	defer cleaner()
	tree, error := core.TraverseDir(root, opt, -1)
	if error != nil {
		t.Errorf("Unable to traverse tree rooted at %s", root)
	}
	if tree.Stats.DirCount != 5 {
		t.Errorf("Expected to get 5 directories but got %d", tree.Stats.DirCount)
	}

	if tree.Stats.FileCount != 6 {
		t.Errorf("Expected to get 6 files but got %d", tree.Stats.FileCount)
	}
}

func TestDirStatIncludingPattern(t *testing.T) {
	cleaner, opt, root := helper.SetupTestDir(uuid.New().String())
	opt.IncludePattern = "*.go"
	defer cleaner()
	tree, error := core.TraverseDir(root, opt, -1)
	if error != nil {
		t.Errorf("Unable to traverse tree rooted at %s", root)
	}
	if tree.Stats.DirCount != 5 {
		t.Errorf("Expected to get 5 directories but got %d", tree.Stats.DirCount)
	}

	if tree.Stats.FileCount != 3 {
		t.Errorf("Expected to get 3 files but got %d", tree.Stats.FileCount)
	}
}

func TestDirStatExcludePattern(t *testing.T) {
	cleaner, opt, root := helper.SetupTestDir(uuid.New().String())
	opt.ExcludePattern = "*.go"
	defer cleaner()
	tree, error := core.TraverseDir(root, opt, -1)
	if error != nil {
		t.Errorf("Unable to traverse tree rooted at %s", root)
	}
	if tree.Stats.DirCount != 5 {
		t.Errorf("Expected to get 5 directories but got %d", tree.Stats.DirCount)
	}

	if tree.Stats.FileCount != 2 {
		t.Errorf("Expected to get 3 files but got %d", tree.Stats.FileCount)
	}
}

func TestDirStatIncludingPatternWithMaxLevel(t *testing.T) {
	cleaner, opt, root := helper.SetupTestDir(uuid.New().String())
	opt.IncludePattern = "*.go"
	opt.MaxLevel = 1
	defer cleaner()
	tree, error := core.TraverseDir(root, opt, -1)
	if error != nil {
		t.Errorf("Unable to traverse tree rooted at %s", root)
	}
	if tree.Stats.DirCount != 3 {
		t.Errorf("Expected to get 3 directories but got %d", tree.Stats.DirCount)
	}

	if tree.Stats.FileCount != 1 {
		t.Errorf("Expected to get 1 files but got %d", tree.Stats.FileCount)
	}
}
