package tree

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

//FilterFunc Signature for filter functions
type filterFunc func(fi os.FileInfo) bool

//Filter utility to filter files
func filter(eles []os.FileInfo, f filterFunc) []os.FileInfo {
	filtered := make([]os.FileInfo, 0)
	for _, e := range eles {
		if f(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func isVisible(fi os.FileInfo) bool {
	return !strings.HasPrefix(path.Base(fi.Name()), ".")
}

func isDir(fi os.FileInfo) bool {
	return fi.IsDir()
}

//FilterInDir return on those file infos which are of type directory
func filterInDir(fis []os.FileInfo) []os.FileInfo {
	return filter(fis, isDir)
}

//FilterOutHidden return on those file which are not hidden
func filterOutHidden(fis []os.FileInfo) []os.FileInfo {
	return filter(fis, isVisible)
}

func filterPattern(fis []os.FileInfo, pattern string, include bool) []os.FileInfo {
	return filter(fis, func(fi os.FileInfo) bool {
		if fi.IsDir() { // Pattern only applies for files
			return true
		}
		matched, err := filepath.Match(pattern, fi.Name())
		if err != nil {
			panic(err)
		}
		if include {
			return matched
		}
		return !matched
	})
}
