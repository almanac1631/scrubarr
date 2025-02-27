package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/folder_scanning"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	"time"
)

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
	statement, err := tx.Prepare("insert into main.entry_mappings (retriever_id, date_added, name, api_resp) values (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("could not prepare entry mappings insert statement: %w", err)
	}
	for retrieverInfo, entries := range rawEntries {
		for name, entry := range entries {
			dateAdded, err := getDateAddedFromEntry(entry)
			var apiResp any
			apiResp, err = json.Marshal(entry.AdditionalData)
			if err != nil {
				return fmt.Errorf("could not marshal entry for retriever (%+v): %w", retrieverInfo, err)
			}
			if _, err = statement.Exec(retrieverInfo.Id(), dateAdded, name, apiResp); err != nil {
				return fmt.Errorf("could not insert entry for retriever (%+v): %w", retrieverInfo, err)
			}
		}
	}
	return nil
}

func getDateAddedFromEntry(entry common.Entry) (*time.Time, error) {
	var dateAdded time.Time
	// switch type of entry.AdditionalData
	switch entry.AdditionalData.(type) {
	case arr_apps.ArrAppEntry:
		dateAdded = entry.AdditionalData.(arr_apps.ArrAppEntry).DateAdded
	case folder_scanning.FileEntry:
		dateAdded = entry.AdditionalData.(folder_scanning.FileEntry).DateModified
	case torrent_clients.TorrentClientEntry:
		dateAdded = entry.AdditionalData.(torrent_clients.TorrentClientEntry).DownloadedAt
	default:
		return nil, fmt.Errorf("could not get added date from entry: unknown entry type %T", entry.AdditionalData)
	}

	if dateAdded.IsZero() {
		return nil, nil
	}
	return &dateAdded, nil
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
