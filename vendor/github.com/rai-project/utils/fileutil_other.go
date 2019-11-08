// +build !linux

package utils

import "os"

func hasHardLinks(fi os.FileInfo) bool {
	return false
}

func getInode(fi os.FileInfo) uint64 {
	return 0
}
