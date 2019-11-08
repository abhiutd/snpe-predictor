// +build linux

package utils

import (
	"os"
	"syscall"
)

func hasHardLinks(fi os.FileInfo) bool {
	// On directories, Nlink doesn't make sense when checking for hard links
	return !fi.IsDir() && fi.Sys().(*syscall.Stat_t).Nlink > 1
}

func getInode(fi os.FileInfo) uint64 {
	return fi.Sys().(*syscall.Stat_t).Ino
}
