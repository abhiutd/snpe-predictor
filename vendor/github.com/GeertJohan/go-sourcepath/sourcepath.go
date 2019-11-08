package sourcepath

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func absoluteDir() (string, error) {
	_, callingGoFile, _, ok := runtime.Caller(2)
	if !ok {
		return "", errors.New("couldn't find caller on stack")
	}

	// resolve to proper path
	pkgDir := filepath.Dir(callingGoFile)
	// fix for go cover
	const coverPath = "_test/_obj_test"
	if !filepath.IsAbs(pkgDir) {
		if i := strings.Index(pkgDir, coverPath); i >= 0 {
			pkgDir = pkgDir[:i] + pkgDir[i+len(coverPath):]            // remove coverPath
			pkgDir = filepath.Join(os.Getenv("GOPATH"), "src", pkgDir) // make absolute
		}
	}
	return pkgDir, nil
}

// MustAbsoluteDir returns the absolute path for the directory in which the calling go file is located.
// When an error occured, this function will panic.
func MustAbsoluteDir() string {
	path, err := absoluteDir()
	if err != nil {
		panic(err)
	}
	return path
}

// AbsoluteDir returns the absolute path for the directory in which the calling go file is located.
// When an error occured, it is returned instead.
func AbsoluteDir() (string, error) {
	return absoluteDir()
}
