package folder_scanning

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFolderScannerClient(t *testing.T) {
	folderScanner, err := NewFolderScanner([]string{".go"}, "/Users/michael/GolandProjects/scrubarr/")
	assert.NoError(t, err)
	entries, err := folderScanner.RetrieveEntries()
	assert.NoError(t, err)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}
