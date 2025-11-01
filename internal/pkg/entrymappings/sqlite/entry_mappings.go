package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/almanac1631/scrubarr/internal/pkg/common"
)

func (e *EntryMappingManager) GetEntryMappings(page int, pageSize int, filter common.EntryMappingFilter, sortBy common.EntryMappingSortBy, name string) (entryMappings []*common.EntryMapping, totalCount int, err error) {
	offset := (page - 1) * pageSize

	categoriesFilter, err := getCategoriesFilter(filter)
	if err != nil {
		return nil, 0, err
	}

	nameFilter := getNameFilter(name)

	sortByAggrColumn, sortByOrderBy, err := getSortBy(sortBy)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`with category_counts as (select em.name,
                                group_concat(distinct r.category order by r.category) as categories%s
                         from entry_mappings em
                                  join retrievers r on em.retriever_id = r.retriever_id%s
                         group by em.name%s),
     filtered_names as (select name from category_counts%s),
     total_count as (select count(distinct name) as total from filtered_names),
     filtered_entries as (select em.*
                          from entry_mappings em
                                   join category_counts cc on em.name = cc.name
                                   join (select name from filtered_names limit ? offset ?) limited_filtered_names
                                        on em.name = limited_filtered_names.name)
select fe.id, fe.name, fe.retriever_id, fe.date_added, fe.size, tc.total
from filtered_entries fe
         join total_count tc%s;`, sortByAggrColumn, nameFilter, sortByOrderBy, categoriesFilter, sortByOrderBy)
	args := []any{
		pageSize, offset,
	}
	if nameFilter != "" {
		args = append([]any{strings.ToLower(name)}, args...)
	}
	return e.queryAndParseEntryMappings(query, args)
}

func (e *EntryMappingManager) queryAndParseEntryMappings(query string, args []any) (entryMappings []*common.EntryMapping, totalCount int, err error) {
	slog.Info("querying entry mappings", "sql", query, "args", args)

	var rows *sql.Rows
	rows, err = e.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("could not query entry mappings: %w", err)
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
		var id, name, retrieverId string
		var dateAdded time.Time
		var size int64
		if err = rows.Scan(&id, &name, &retrieverId, &dateAdded, &size, &totalCount); err != nil {
			entryMappings = nil
			err = fmt.Errorf("could not scan entry mappings: %w", err)
			return
		}
		entryMappings, err = e.parseEntryMapping(id, name, retrieverId, dateAdded, size, entryMappings)
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

func getSortBy(sortBy common.EntryMappingSortBy) (string, string, error) {
	var sortAscending bool
	var sortByAggr, sortByColName string
	switch sortBy {
	case common.EntryMappingSortByNoSort:
		return "", "", nil
	case common.EntryMappingSortByDateAsc, common.EntryMappingSortByDateDesc:
		sortByAggr = "max(em.date_added)"
		sortByColName = "date_added"
		sortAscending = sortBy == common.EntryMappingSortByDateAsc
		break
	case common.EntryMappingSortBySizeAsc, common.EntryMappingSortBySizeDesc:
		sortByAggr = "max(em.size)"
		sortByColName = "size"
		sortAscending = sortBy == common.EntryMappingSortBySizeAsc
		break
	case common.EntryMappingSortByNameAsc, common.EntryMappingSortByNameDesc:
		sortByAggr = "em.name"
		sortByColName = "name"
		sortAscending = sortBy == common.EntryMappingSortByNameAsc
	default:
		return "", "", fmt.Errorf("invalid sort by: %d", sortBy)
	}
	orderBySuffix := "desc"
	if sortAscending {
		orderBySuffix = "asc"
	}
	return fmt.Sprintf(", %s as %s", sortByAggr, sortByColName),
		fmt.Sprintf(" order by %s %s", sortByColName, orderBySuffix), nil
}

func getCategoriesFilter(filter common.EntryMappingFilter) (string, error) {
	categoriesFilter := " where categories %s (select group_concat(distinct r.category order by r.category) as categories from retrievers r)"
	if filter == common.EntryMappingFilterNoFilter {
		categoriesFilter = ""
	} else if filter == common.EntryMappingFilterCompleteEntry {
		categoriesFilter = fmt.Sprintf(categoriesFilter, "=")
	} else if filter == common.EntryMappingFilterIncompleteEntry {
		categoriesFilter = fmt.Sprintf(categoriesFilter, "!=")
	} else {
		return "", fmt.Errorf("invalid filter: %d", filter)
	}
	return categoriesFilter, nil
}

func getNameFilter(name string) string {
	if name == "" {
		return ""
	}
	return " where LOWER(em.name) LIKE ('%' || LOWER(replace(REPLACE(?, ' ', '%'), '_', '\\_')) || '%') ESCAPE '\\'"
}

func (e *EntryMappingManager) GetEntryMappingById(id string) (*common.EntryMapping, error) {
	query := `select em.id, em.name, em.retriever_id, em.date_added, em.size, 0 from entry_mappings em where em.id = ?;`
	args := []any{id}
	entryMappings, _, err := e.queryAndParseEntryMappings(query, args)
	if err != nil {
		return nil, fmt.Errorf("could not get entry mapping by id %q: %w", id, err)
	}
	if len(entryMappings) == 0 {
		return nil, common.ErrEntryMappingNotFound
	}
	return entryMappings[0], nil
}

func (e *EntryMappingManager) GetRetrieverById(id common.RetrieverId) (common.RetrieverInfo, common.EntryRetriever, error) {
	for retrieverInfo, retriever := range e.entryRetrievers {
		if retrieverInfo.Id() != id {
			continue
		}
		return retrieverInfo, retriever, nil
	}
	return common.RetrieverInfo{}, nil, errors.New("retriever not found")
}

func (e *EntryMappingManager) parseEntryMapping(id string, entryName string, retrieverId string, dateAdded time.Time, size int64, entryMappings []*common.EntryMapping) ([]*common.EntryMapping, error) {
	retrieverIdMapped := common.RetrieverId(retrieverId)

	retrieverInfo, _, err := e.GetRetrieverById(retrieverIdMapped)
	if err != nil {
		return nil, err
	}
	var entryMapping *common.EntryMapping
	for _, presentEntryMapping := range entryMappings {
		if presentEntryMapping.Id == id {
			entryMapping = presentEntryMapping
			break
		}
	}
	if entryMapping == nil {
		entryMapping = &common.EntryMapping{
			Id:              id,
			Name:            common.EntryName(entryName),
			RetrieversFound: []common.RetrieverInfo{},
		}
		entryMappings = append(entryMappings, entryMapping)
	}
	entryMapping.RetrieversFound = append(entryMapping.RetrieversFound, retrieverInfo)
	if entryMapping.DateAdded.IsZero() || dateAdded.Before(entryMapping.DateAdded) {
		entryMapping.DateAdded = dateAdded
	}
	if entryMapping.Size == 0 || size > entryMapping.Size {
		entryMapping.Size = size
	}
	return entryMappings, nil
}

func (e *EntryMappingManager) GetEntryMappingDetails(id string) (details common.EntryMappingDetails, err error) {
	query := `select em.retriever_id, em.api_resp from entry_mappings em where em.id = ?;`
	var rows *sql.Rows
	rows, err = e.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("could not get entry mapping details by id %q: %w", id, err)
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if err != nil {
			err = closeErr
		}
	}(rows)
	details = make(common.EntryMappingDetails)
	for rows.Next() {
		if rows.Err() != nil {
			return nil, fmt.Errorf("could not get entry mapping details by id %q: %w", id, rows.Err())
		}
		var retrieverId, apiResp string
		if err = rows.Scan(&retrieverId, &apiResp); err != nil {
			return nil, fmt.Errorf("could not scan entry mapping details: %w", err)
		}
		details[common.RetrieverId(retrieverId)] = apiResp
	}
	if len(details) == 0 {
		return nil, common.ErrEntryMappingNotFound
	}
	return details, nil
}
