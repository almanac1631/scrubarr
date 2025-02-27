package sqlite

import (
	"database/sql"
	"embed"
	"encoding/json"
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

func (e *EntryMappingManager) GetEntryMappings(page int, pageSize int, filter common.EntryMappingFilter) (entryMappings []*common.EntryMapping, totalCount int, err error) {
	offset := (page - 1) * pageSize

	var rows *sql.Rows
	categoriesFilter := " where categories %s (select group_concat(distinct r.category order by r.category) as categories from retrievers r)"
	if filter == common.EntryMappingFilterNoFilter {
		categoriesFilter = ""
	} else if filter == common.EntryMappingFilterCompleteEntry {
		categoriesFilter = fmt.Sprintf(categoriesFilter, "=")
	} else if filter == common.EntryMappingFilterIncompleteEntry {
		categoriesFilter = fmt.Sprintf(categoriesFilter, "!=")
	} else {
		return nil, 0, fmt.Errorf("invalid filter: %d", filter)
	}

	query := fmt.Sprintf(`with category_counts as (select em.name, group_concat(distinct r.category order by r.category) as categories
                         from entry_mappings em
                                  join retrievers r on em.retriever_id = r.retriever_id
                         group by em.name order by em.name),
     filtered_names as (select name from category_counts%s),
     total_count as (select count(distinct name) as total from filtered_names),
     filtered_entries as (select em.*
                          from entry_mappings em
                                   join category_counts cc on em.name = cc.name
                                   join (select name from filtered_names limit ? offset ?) limited_filtered_names
                                        on em.name = limited_filtered_names.name)
select fe.name, fe.retriever_id, tc.total
from filtered_entries fe
         join total_count tc order by fe.name;`, categoriesFilter)
	rows, err = e.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("could not get entry mappings: %w", err)
	}

	defer func() {
		if err != nil {
			_ = rows.Close()
		} else {
			err = rows.Close()
		}
	}()

	entryMappings = []*common.EntryMapping{}
	for rows.Next() {
		var name, retrieverId string
		if err = rows.Scan(&name, &retrieverId, &totalCount); err != nil {
			entryMappings = nil
			err = fmt.Errorf("could not scan entry mappings: %w", err)
			return
		}
		entryMappings, err = e.parseEntryMapping(name, retrieverId, entryMappings)
		if err != nil {
			entryMappings = nil
			err = fmt.Errorf("could not parse entry mappings: %w", err)
			return
		}
	}
	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("could not get entry mappings: %w", err)
	}
	return
}

func (e *EntryMappingManager) parseEntryMapping(entryName string, retrieverId string, entryMappings []*common.EntryMapping) ([]*common.EntryMapping, error) {
	entryNameMapped := common.EntryName(entryName)
	retrieverIdMapped := common.RetrieverId(retrieverId)

	retrieverInfo, err := e.getRetrieverById(retrieverIdMapped)
	if err != nil {
		return nil, err
	}
	var entryMapping *common.EntryMapping
	for _, presentEntryMapping := range entryMappings {
		if presentEntryMapping.Name == entryNameMapped {
			entryMapping = presentEntryMapping
			break
		}
	}
	if entryMapping == nil {
		entryMapping = &common.EntryMapping{
			Name:            entryNameMapped,
			RetrieversFound: []common.RetrieverInfo{},
		}
		entryMappings = append(entryMappings, entryMapping)
	}
	entryMapping.RetrieversFound = append(entryMapping.RetrieversFound, retrieverInfo)
	return entryMappings, nil
}

func (e *EntryMappingManager) getRetrieverById(retrieverId common.RetrieverId) (common.RetrieverInfo, error) {
	for retrieverInfo, _ := range e.entryRetrievers {
		if retrieverInfo.Id() == retrieverId {
			return retrieverInfo, nil
		}
	}
	return common.RetrieverInfo{}, fmt.Errorf("retriever with id %s not found", retrieverId)
}

func (e *EntryMappingManager) GetRetrievers() ([]common.RetrieverInfo, error) {
	return slices.Collect(maps.Keys(e.entryRetrievers)), nil
}
