package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	_ "github.com/glebarez/go-sqlite"
	"github.com/pressly/goose/v3"
	"maps"
	"slices"
)

var _ common.EntryMappingManager = (*EntryMappingManager)(nil)

type EntryMappingManager struct {
	entryRetrievers       map[common.RetrieverInfo]common.EntryRetriever
	bundledEntryRetriever common.BundledEntryRetriever

	db *sql.DB
}

func NewEntryMappingManager(entryRetrievers map[common.RetrieverInfo]common.EntryRetriever, bundledEntryRetriever common.BundledEntryRetriever, sqlitePath string) (*EntryMappingManager, error) {
	manager := &EntryMappingManager{
		entryRetrievers:       entryRetrievers,
		bundledEntryRetriever: bundledEntryRetriever,
	}

	var err error
	manager.db, err = sql.Open("sqlite", sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}
	if err := manager.applyMigrations(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return manager, nil
}

func (e *EntryMappingManager) Close() error {
	if e.db != nil {
		return e.db.Close()
	}
	return nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func (e *EntryMappingManager) applyMigrations() error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.Up(e.db, "migrations"); err != nil {
		return err
	}
	return nil
}

func (e *EntryMappingManager) GetRetrievers() ([]common.RetrieverInfo, error) {
	return slices.Collect(maps.Keys(e.entryRetrievers)), nil
}
