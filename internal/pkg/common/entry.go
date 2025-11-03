package common

type Entry struct {
	Name           EntryName
	FilePath       string
	ParentId       string
	AdditionalData any
}

type EntryName string
