//go:build windows

package utils

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetFileInode(path string) (uint64, error) {
	utf16Path, err := windows.UTF16FromString(path)
	if err != nil {
		return 0, fmt.Errorf("failed to convert path to utf-16: %v", err)
	}

	handle, err := windows.CreateFile(
		&utf16Path[0],
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to open file handle: %v", err)
	}
	defer windows.Close(handle)

	var fileInfo windows.ByHandleFileInformation
	err = windows.GetFileInformationByHandleEx(handle, windows.FileIdInfo, (*byte)(unsafe.Pointer(&fileInfo)), uint32(unsafe.Sizeof(fileInfo)))
	if err != nil {
		return 0, fmt.Errorf("failed to get file information: %v", err)
	}

	fileID := uint64(fileInfo.FileIndexHigh)<<32 | uint64(fileInfo.FileIndexLow)

	return fileID, nil
}
