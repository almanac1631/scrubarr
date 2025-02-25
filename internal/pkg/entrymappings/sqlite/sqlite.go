package sqlite

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	_ "github.com/glebarez/go-sqlite"
	"github.com/pressly/goose/v3"
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

func (e *EntryMappingManager) RefreshEntryMappings() (err error) {
	var rawEntries map[common.RetrieverInfo]common.RetrieverEntries
	rawEntries, err = e.bundledEntryRetriever(e.entryRetrievers)
	if err != nil {
		err = fmt.Errorf("could not query entries using given entry retriever: %w", err)
		return
	}
	var tx *sql.Tx
	tx, err = e.db.Begin()
	defer func() {
		if err == nil {
			return
		}
		if err = tx.Rollback(); err != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", err)
		}
	}()

	err = e.updateRetrievers(tx)
	if err != nil {
		return
	}

	err = e.updateEntryMappings(tx, rawEntries)
	if err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("could not commit transaction: %w", err)
		return
	}
	return nil
}

func (e *EntryMappingManager) updateEntryMappings(tx *sql.Tx, rawEntries map[common.RetrieverInfo]common.RetrieverEntries) error {
	//goland:noinspection SqlWithoutWhere
	if _, err := tx.Exec("delete from main.entry_mappings;"); err != nil {
		return fmt.Errorf("could not truncate entry_mappings table: %w", err)
	}
	statement, err := tx.Prepare("insert into main.entry_mappings (retriever_id, name, api_resp) values (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare entry mappings insert statement: %w", err)
	}
	for retrieverInfo, entries := range rawEntries {
		for name, entry := range entries {
			var apiResp any
			apiResp, err = json.Marshal(entry.AdditionalData)
			if err != nil {
				return fmt.Errorf("could not marshal entry for retriever (%+v): %w", retrieverInfo, err)
			}
			if _, err = statement.Exec(retrieverInfo.Id(), name, apiResp); err != nil {
				return fmt.Errorf("could not insert entry for retriever (%+v): %w", retrieverInfo, err)
			}
		}
	}
	return nil
}

func (e *EntryMappingManager) updateRetrievers(tx *sql.Tx) error {
	//goland:noinspection SqlWithoutWhere
	_, err := tx.Exec("delete from main.retrievers;")
	if err != nil {
		return fmt.Errorf("could not truncate retrievers table: %w", err)
	}
	statement, err := tx.Prepare("insert into main.retrievers (retriever_id, category, software_name, name) values (?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("could not prepare retriever insert statement: %w", err)
	}
	for retriever := range e.entryRetrievers {
		if _, err = statement.Exec(retriever.Id(), retriever.Category, retriever.SoftwareName, retriever.Name); err != nil {
			return fmt.Errorf("could not insert retriever (%+v): %w", retriever, err)
		}
	}
	return nil
}

func (e *EntryMappingManager) GetEntryMappings(page int, pageSize int, filter common.EntryMappingFilter) ([]common.EntryMapping, int, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EntryMappingManager) GetRetrievers() ([]common.RetrieverInfo, error) {
	//TODO implement me
	panic("implement me")
}
