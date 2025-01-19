package config

import (
	"fmt"
)

type EntryNotFound struct {
	path string
}

func (e EntryNotFound) Error() string {
	return fmt.Sprintf("config entry not found (path: %q)", e.path)
}

type EntryAccessor func(path string) interface{}

func GetEntry[V any](accessor EntryAccessor, path string) (mappedValue V, err error) {
	rawValue := accessor(path)
	if rawValue == nil {
		return mappedValue, &EntryNotFound{path: path}
	}
	var ok bool
	mappedValue, ok = rawValue.(V)
	if !ok {
		return mappedValue, fmt.Errorf("could not apply config entry type (path: %q, type: %T)", path, *new(V))
	}
	return mappedValue, nil
}
