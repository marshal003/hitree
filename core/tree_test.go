package core_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/marshal003/hitree/core"
	"github.com/marshal003/hitree/core/helper"
)

func TestTreePrint(t *testing.T) {
	root_file := uuid.New().String()
	cleaner, opt, root := helper.SetupTestDir(root_file)
	defer cleaner()
	tree, err := core.TraverseDir(root, opt, -1)
	if err != nil {
		t.Errorf("Unable to traverse the Dir")
	}
	buf := new(bytes.Buffer)
	tree.Print(buf, opt)
	output := []byte(fmt.Sprintf(`%s
├──a
│  ├──b
│  │  └──normal.go
│  ├──c
│  │  ├──d
│  │  │  ├──e
│  │  │  └──normal.py
│  │  └──normal.go
│  └──normal.py
└──normal.go

5 directories, 5 files
`, root_file))
	if !bytes.Equal(output, buf.Bytes()) {
		t.Errorf("expected %s, got: \n%s", output, buf.Bytes())
	}
}

func TestTreePrint_prune(t *testing.T) {
	root_file := uuid.New().String()
	cleaner, opt, root := helper.SetupTestDir(root_file)
	defer cleaner()
	opt.Prune = true
	tree, err := core.TraverseDir(root, opt, -1)
	if err != nil {
		t.Errorf("Unable to traverse the Dir")
	}
	buf := new(bytes.Buffer)
	tree.Print(buf, opt)
	output := []byte(fmt.Sprintf(`%s
├──a
│  ├──b
│  │  └──normal.go
│  ├──c
│  │  ├──d
│  │  │  └──normal.py
│  │  └──normal.go
│  └──normal.py
└──normal.go

5 directories, 5 files
`, root_file))
	if !bytes.Equal(output, buf.Bytes()) {
		t.Errorf("expected %s, got: \n%s", output, buf.Bytes())
	}
}

func TestTreePrint_level(t *testing.T) {
	root_file := uuid.New().String()
	cleaner, opt, root := helper.SetupTestDir(root_file)
	defer cleaner()
	opt.MaxLevel = 1
	tree, err := core.TraverseDir(root, opt, -1)
	if err != nil {
		t.Errorf("Unable to traverse the Dir")
	}
	buf := new(bytes.Buffer)
	tree.Print(buf, opt)
	output := []byte(fmt.Sprintf(`%s
├──a
│  ├──b
│  ├──c
│  └──normal.py
└──normal.go

3 directories, 2 files
`, root_file))
	if !bytes.Equal(output, buf.Bytes()) {
		t.Errorf("expected %s, got: \n%s", output, buf.Bytes())
	}
}

func TestJsonTree(t *testing.T) {
	root_file := uuid.New().String()
	cleaner, opt, root := helper.SetupTestDir(root_file)
	defer cleaner()
	opt.DirOnly = true
	tree, _ := core.TraverseDir(root, opt, -1)
	jsonTree := tree.AsJSONTree(opt)
	expectedJSON := fmt.Sprintf(`
		{
			"name": "%s",
			"file_type": "dir",
			"subtree": [
				{
					"name": "a",
					"file_type": "dir",
					"subtree": [
						{
							"name": "b",
							"file_type": "dir",
							"subtree": [
							]
						},
						{
							"name": "c",
							"file_type": "dir",
							"subtree": [
								{
									"name": "d",
									"file_type": "dir",
									"subtree": [
										{
											"name": "e",
											"file_type": "dir",
											"subtree": [
											]
										}
									]
								}
							]
						}
					]
				}
			]
		}`, root_file)
	var expectedJSONTree core.JSONTree
	if err := json.Unmarshal([]byte(expectedJSON), &expectedJSONTree); err != nil {
		t.Errorf("Error: %s", err)
	}
	if !reflect.DeepEqual(expectedJSONTree, jsonTree) {
		t.Errorf("Expected %v, got %v", expectedJSONTree, jsonTree)
	}
}

func TestJsonTree_prune(t *testing.T) {
	root_file := uuid.New().String()
	cleaner, opt, root := helper.SetupTestDir(root_file)
	defer cleaner()
	opt.DirOnly = true
	opt.Prune = true
	tree, _ := core.TraverseDir(root, opt, -1)
	jsonTree := tree.AsJSONTree(opt)
	expectedJSON := fmt.Sprintf(`
		{
			"name": "%s",
			"file_type": "dir",
			"subtree": [
				{
					"name": "a",
					"file_type": "dir",
					"subtree": [
						{
							"name": "b",
							"file_type": "dir",
							"subtree": [
							]
						},
						{
							"name": "c",
							"file_type": "dir",
							"subtree": [
								{
									"name": "d",
									"file_type": "dir",
									"subtree": [
									]
								}
							]
						}
					]
				}
			]
		}`, root_file)
	var expectedJSONTree core.JSONTree
	if err := json.Unmarshal([]byte(expectedJSON), &expectedJSONTree); err != nil {
		t.Errorf("Error: %s", err)
	}
	if !reflect.DeepEqual(expectedJSONTree, jsonTree) {
		t.Errorf("Expected %v, got %v", expectedJSONTree, jsonTree)
	}
}
