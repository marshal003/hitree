package tree

import (
	"io/ioutil"
	"os"
	"path"
)

//TraverseDir utility method to recursively traverse through the dir
func TraverseDir(root string, opt Options, level int16) (Tree, error) {
	var tree Tree
	fi, err := fileStat(root, opt)
	if err != nil {
		return tree, err
	}
	if !fi.IsDir() {
		return Tree{Root: fi, Stats: Stats{0, 0}}, nil
	}
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return tree, err
	}

	files = applyFilters(files, opt)
	childrens := make([]Tree, 0)
	stats := Stats{}

	for _, fi := range files {
		if opt.MaxLevel > -1 && level >= opt.MaxLevel {
			continue
		}
		//DFS of tree
		tree, err := TraverseDir(path.Join(root, fi.Name()), opt, level+1)
		if err != nil {
			return tree, err
		}
		stats = updateStats(tree, stats)
		childrens = updateChildrens(tree, childrens, opt, fi)
	}
	tree = Tree{Root: fi, Childrens: childrens, Stats: stats}
	return tree, nil
}

func fileStat(path string, opt Options) (os.FileInfo, error) {
	if opt.FollowLink {
		return os.Stat(path)
	}
	return os.Lstat(path)
}

func updateChildrens(tree Tree, childrens []Tree, opt Options, fi os.FileInfo) []Tree {
	if opt.DirOnly && fi.IsDir() {
		childrens = append(childrens, tree)
	} else if !opt.DirOnly {
		childrens = append(childrens, tree)
	}
	return childrens
}

func updateStats(tree Tree, stats Stats) Stats {
	stats.DirCount = stats.DirCount + tree.Stats.DirCount
	stats.FileCount = stats.FileCount + tree.Stats.FileCount
	if tree.Root.IsDir() {
		stats.DirCount++
	} else {
		stats.FileCount++
	}
	return stats
}

func applyFilters(fis []os.FileInfo, opt Options) []os.FileInfo {
	if !opt.IncludeHidden {
		fis = filterOutHidden(fis)
	}

	if len(opt.ExcludePattern) > 0 {
		fis = filterPattern(fis, opt.ExcludePattern, false)
	}

	if len(opt.IncludePattern) > 0 {
		fis = filterPattern(fis, opt.IncludePattern, true)
	}
	return fis
}
