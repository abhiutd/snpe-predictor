package utils

import (
	"os"
	"path/filepath"
)

func FileSize(filepath string) (int64, error) {
	file, err := os.Open(filepath)
	defer file.Close()

	if err != nil {
		return int64(0), err
	}

	fi, err := file.Stat()
	if err != nil {
		return int64(0), err
	}

	return fi.Size(), nil
}

// DirSize takes a path and returns its size in bytes
func DirSize(path string) (int64, error) {
	seenInode := make(map[uint64]struct{})

	if _, err := os.Stat(path); err == nil {
		var sz int64
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if hasHardLinks(info) {
				ino := getInode(info)
				if _, ok := seenInode[ino]; !ok {
					seenInode[ino] = struct{}{}
					sz += info.Size()
				}
			} else {
				sz += info.Size()
			}
			return err
		})
		return sz, err
	}

	return 0, nil
}
