package common

type Entry struct {
	Name           EntryName
	FilePath       string
	AdditionalData any
}

type EntryName string
