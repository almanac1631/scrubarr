//go:build unix

package utils

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func GetFileInode(path string) (uint64, error) {
	var stat unix.Stat_t
	err := unix.Stat(path, &stat)
	if err != nil {
		return 0, fmt.Errorf("could not stat file %s: %v", path, err)
	}

	return stat.Ino, nil
}
