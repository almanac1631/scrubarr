package folder_scanning

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"os"
	pathLib "path"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type FileEntry struct {
	Path        string
	SizeInBytes int64
	DateAdded   time.Time
}

type FolderScanner struct {
	folderPath         string
	allowedFileEndings []string
}

func NewFolderScanner(allowedFileEndings []string, folderPath string) (*FolderScanner, error) {
	folderScanner := &FolderScanner{folderPath, allowedFileEndings}
	fileInfo, err := os.Stat(folderPath)
	if err != nil {
		return nil, fmt.Errorf("could not stat folder path: %w", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("given folder path is no directory: %q", folderPath)
	}
	return folderScanner, nil
}

func (folderScanner *FolderScanner) RetrieveEntries() (common.RetrieverEntries, error) {
	entries := common.RetrieverEntries{}
	err := filepath.Walk(folderScanner.folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ext := pathLib.Ext(path)
		if !slices.Contains(folderScanner.allowedFileEndings, strings.ToLower(ext)) {
			return nil
		}
		name := pathLib.Base(path)
		nameLower := strings.ToLower(name)
		entry := common.Entry{
			Name: common.EntryName(name),
			AdditionalData: FileEntry{
				Path:        path,
				SizeInBytes: info.Size(),
				DateAdded:   info.ModTime(),
			},
		}
		entries[common.EntryName(nameLower)] = entry
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}
