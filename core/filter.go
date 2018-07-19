package core

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

//FilterFunc Signature for filter functions
type filterFunc func(fi os.FileInfo) bool

//FileFilterFunc ...
type FileFilterFunc func(fis []os.FileInfo) []os.FileInfo

//Filter utility to Filter files
func Filter(eles []os.FileInfo, f filterFunc) []os.FileInfo {
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
func FilterInDir(fis []os.FileInfo) []os.FileInfo {
	return Filter(fis, isDir)
}

//FilterOutHidden return on those file which are not hidden
func FilterOutHidden(fis []os.FileInfo) []os.FileInfo {
	return Filter(fis, isVisible)
}

// FileFilterPattern Filter fileInfos whose filepath matched with patern. Include flag
// is used to control the inclusion/exclusion of matched fileinfos in the result set.
// eg. If include flag is true, and file path matched with given pattern, then that
// file info will be incuded in the result.
// Note: This filter will only be applied on files
func FileFilterPattern(fis []os.FileInfo, pattern string, include bool) []os.FileInfo {
	return Filter(fis, func(fi os.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		matched, err := filepath.Match(pattern, fi.Name())
		if err != nil {
			panic("\nRegularExpression is incorrect, Please refer https://golang.org/pkg/path/filepath/#Match\n")
		}
		if include {
			return matched
		}
		return !matched
	})
}
