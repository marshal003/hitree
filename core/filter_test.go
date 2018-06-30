package core_test

import (
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	"github.com/marshal003/hitree/core"
	"github.com/marshal003/hitree/core/helper"
)

// SetupDir returns directory structure
//
// root
// 	.hidden
// 	golang.go
// 	python.py
// 	python.pyc
//  normal.py
// 	normal.go
//  normal.txt
//  a -> dir
func TestFilterIndir(t *testing.T) {
	cleaner, _, root := helper.SetupFlattenTestDir(uuid.New().String())
	defer cleaner()
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		t.Errorf("Unable to read files from %s, %v", root, err)
	}
	if len(fis) != 8 {
		t.Errorf("Expected files count to %d, got %d", 8, len(fis))
	}
	filterFis := core.FilterInDir(fis)
	if len(filterFis) != 1 {
		t.Errorf("There should have been only one directory, but got %d", len(filterFis))
	}
}

func TestFilterOutHidden(t *testing.T) {
	cleaner, _, root := helper.SetupFlattenTestDir(uuid.New().String())
	defer cleaner()
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		t.Errorf("Unable to read files from %s, %v", root, err)
	}
	if len(fis) != 8 {
		t.Errorf("Expected files count to %d, got %d", 8, len(fis))
	}
	filterFis := core.FilterOutHidden(fis)
	if len(filterFis) != 7 {
		t.Errorf("There should have been only %d file infos, but got %d", 7, len(filterFis))
	}
}

func TestFileFilterPattern(t *testing.T) {
	cleaner, _, root := helper.SetupFlattenTestDir(uuid.New().String())
	defer cleaner()
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		t.Errorf("Unable to read files from %s, %v", root, err)
	}
	if len(fis) != 8 {
		t.Errorf("Expected files count to %d, got %d", 8, len(fis))
	}

	table := []struct {
		expeted int
		pattern string
		include bool
	}{
		{3, "*.go", true},
		{6, "*.go", false},
		{2, "*.txt", true},
		{7, "*.txt", false},
		{2, "*.pyc", true},
		{7, "*.pyc", false},
		{2, "gol*", true},
		{7, "gol*", false},
	}

	for _, test := range table {
		filterFis := core.FileFilterPattern(fis, test.pattern, test.include) // Directory a will not be filtered out
		if len(filterFis) != test.expeted {
			t.Errorf("There should have been only %d file infos, but got %d", test.expeted, len(filterFis))
		}
	}
}
