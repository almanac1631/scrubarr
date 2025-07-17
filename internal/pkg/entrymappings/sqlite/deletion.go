package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
)

func (e *EntryMappingManager) DeleteEntryMappingById(id string) error {
	rows, err := e.db.Query("select retriever_id, api_resp from entry_mappings where id = ?;", id)
	if err != nil {
		return fmt.Errorf("could not query entry mapping: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("could not close rows: %v\n", err)
		}
	}(rows)
	entryFound := false
	for rows.Next() {
		if rows.Err() != nil {
			return fmt.Errorf("could not iterate entry mapping: %w", rows.Err())
		}
		entryFound = true
		var retrieverId, apiResp string
		if err = rows.Scan(&retrieverId, &apiResp); err != nil {
			return fmt.Errorf("could not scan entry mapping: %w", err)
		}
		retrieverInfo, retriever, err := e.GetRetrieverById(common.RetrieverId(retrieverId))
		if err != nil {
			return fmt.Errorf("could not get retriever by id: %w", err)
		}
		retrieverEntryId, err := getIdFromApiResp(retrieverInfo, apiResp)
		if err != nil {
			return fmt.Errorf("could not get id from api response: %w", err)
		}
		err = retriever.DeleteEntry(retrieverEntryId)
		if err != nil {
			return fmt.Errorf("could not delete entry from retriever: %w", err)
		}
	}
	if !entryFound {
		return common.ErrEntryNotFound
	}
	_, err = e.db.Exec("delete from entry_mappings where id = ?;", id)
	if err != nil {
		return fmt.Errorf("could not delete entry mapping from db: %w", err)
	}
	return nil
}

func getIdFromApiResp(retrieverInfo common.RetrieverInfo, apiResp string) (any, error) {
	switch retrieverInfo.Category {
	case "arr_app":
		var arrAppEntry arr_apps.ArrAppEntry
		err := json.Unmarshal([]byte(apiResp), &arrAppEntry)
		return arrAppEntry.ID, err
	case "torrent_client":
		var torrentClientEntry torrent_clients.TorrentClientEntry
		err := json.Unmarshal([]byte(apiResp), &torrentClientEntry)
		return torrentClientEntry.ID, err
	default:
		return nil, fmt.Errorf("unknown retriever category %q", retrieverInfo.Category)
	}
}
