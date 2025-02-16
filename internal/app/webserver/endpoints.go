package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/folder_scanning"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	"log/slog"
	"maps"
	"slices"
	"sort"
	"strings"
)

func (a ApiEndpointHandler) GetEntryMappings(ctx context.Context, request GetEntryMappingsRequestObject) (resp GetEntryMappingsResponseObject, err error) {
	page := request.Params.Page
	pageSize := request.Params.PageSize
	if err = validatePage(request.Params.Page); err != nil {
		return resp, err
	}
	if err = validatePageSize(10); err != nil {
		return resp, err
	}
	entryMappingList := a.retrieverRegistry.RetrieveEntryMapping()
	responseEntryMappingList := make([]EntryMapping, 0)
	entryNameList := slices.Collect(maps.Keys(entryMappingList))
	sort.SliceStable(entryNameList, func(i, j int) bool {
		return strings.Compare(string(entryNameList[i]), string(entryNameList[j])) == -1
	})
	for _, entryName := range entryNameList {
		presencePairs := entryMappingList[entryName]
		entryMapping := getResponseEntryMappingFromPresencePairs(entryName, presencePairs)
		responseEntryMappingList = append(responseEntryMappingList, entryMapping)
	}
	responseEntryMappingList = a.applyFilter(responseEntryMappingList, request.Params.Filter)
	totalAmount := len(responseEntryMappingList)
	begin := page * pageSize
	end := (page + 1) * pageSize
	var responseEntryMappingListPortion []EntryMapping

	if begin > len(responseEntryMappingList) {
		begin = len(responseEntryMappingList)
	}
	if end > len(responseEntryMappingList) {
		end = len(responseEntryMappingList)
	}
	responseEntryMappingListPortion = responseEntryMappingList[begin:end]
	return &GetEntryMappings200JSONResponse{
		responseEntryMappingListPortion,
		totalAmount,
	}, nil
}

func (a ApiEndpointHandler) applyFilter(responseEntryMappingList []EntryMapping, filter *GetEntryMappingsParamsFilter) []EntryMapping {
	if filter == nil {
		return responseEntryMappingList
	}
	filteredEntries := make([]EntryMapping, 0)
	for _, entryMapping := range responseEntryMappingList {
		entryComplete := a.isEntryMappingComplete(entryMapping)
		entryIncluded := (*filter == CompleteEntries && entryComplete) || (*filter == IncompleteEntries && !entryComplete)
		if entryIncluded {
			filteredEntries = append(filteredEntries, entryMapping)
		}
	}
	return filteredEntries
}

func (a ApiEndpointHandler) isEntryMappingComplete(entryMapping EntryMapping) bool {
	categoryFulfilledMap := make(map[string]bool)
	for _, retrieverInfo := range a.retrieverRegistry.GetRetrievers() {
		if _, ok := categoryFulfilledMap[retrieverInfo.Category]; !ok {
			categoryFulfilledMap[retrieverInfo.Category] = false
		}
		for _, retrieverFinding := range entryMapping.RetrieverFindings {
			if common.RetrieverId(retrieverFinding.Id) == retrieverInfo.Id() {
				categoryFulfilledMap[retrieverInfo.Category] = true
			}
		}
	}
	for _, categoryFulfilled := range categoryFulfilledMap {
		if !categoryFulfilled {
			return false
		}
	}
	return true
}

func validatePage(page int) error {
	if page < 0 {
		return &InvalidParamFormatError{
			ParamName: "page",
			Err:       errors.New("page param has to be greater than 0"),
		}
	}
	return nil
}

func validatePageSize(pageSize int) error {
	if pageSize > 100 || pageSize < 10 {
		return &InvalidParamFormatError{
			ParamName: "pageSize",
			Err:       errors.New("page param has to be between 10 and 100"),
		}
	}
	return nil
}

func getResponseEntryMappingFromPresencePairs(entryName retrieval.EntryName, pairs common.EntryPresencePairs) EntryMapping {
	findings := make([]EntryMappingRetrieverFindingsInner, 0)
	for retrieverInfo, entry := range pairs {
		retrieverId := RetrieverId(retrieverInfo.Id())
		mappedFindingValue, err := getFindingValueFromEntry(entry)
		if err != nil {
			slog.Error("failed automatic finding value mapping", "err", err)
			continue
		}
		marshalledFindingValue, err := json.Marshal(mappedFindingValue)
		if err != nil {
			slog.Error("could not marshal finding value", "err", err)
			continue
		}
		finding := EntryMappingRetrieverFindingsInner{
			Id: retrieverId,
			Detail: EntryMappingRetrieverFindingsInnerDetail{
				union: marshalledFindingValue,
			},
		}
		findings = append(findings, finding)
	}
	entryNameStr := string(entryName)
	return EntryMapping{
		Name:              entryNameStr,
		RetrieverFindings: findings,
	}
}

func getFindingValueFromEntry(entry *retrieval.Entry) (any, error) {
	switch entryMapped := entry.AdditionalData.(type) {
	case arr_apps.ArrAppEntry:
		mediaType := ArrAppFindingMediaType(entryMapped.Type)
		return ArrAppFinding{
			MediaFilePath: &entryMapped.MediaFilePath,
			MediaType:     &mediaType,
			Monitored:     &entryMapped.Monitored,
			ParentName:    &entryMapped.ParentName,
		}, nil
	case folder_scanning.FileEntry:
		return FolderFinding{
			FilePath: &entryMapped.Path,
			Size:     &entryMapped.SizeInBytes,
		}, nil
	case torrent_clients.TorrentClientEntry:
		return TorrentClientFinding{
			ClientName:       &entryMapped.TorrentClientName,
			DownloadFilePath: &entryMapped.DownloadFilePath,
			DownloadedAt:     &entryMapped.DownloadedAt,
			Ratio:            &entryMapped.Ratio,
			Size:             &entryMapped.FileSizeBytes,
			TorrentName:      &entryMapped.TorrentName,
			TrackerHost:      &entryMapped.TrackerHost,
		}, nil
	default:
		return nil, fmt.Errorf("could not map unknown entry info type: %T", entry.AdditionalData)
	}
}

func (a ApiEndpointHandler) GetRetrievers(ctx context.Context, request GetRetrieversRequestObject) (GetRetrieversResponseObject, error) {
	var retrieverList []Retriever
	for _, retrieverInfo := range a.retrieverRegistry.GetRetrievers() {
		id := RetrieverId(retrieverInfo.Id())
		category := RetrieverCategory(retrieverInfo.Category)
		name := retrieverInfo.Name
		softwareName := RetrieverSoftwareName(retrieverInfo.SoftwareName)
		retrieverInstance := Retriever{
			Id:           id,
			Category:     category,
			Name:         name,
			SoftwareName: softwareName,
		}
		retrieverList = append(retrieverList, retrieverInstance)
	}
	return GetRetrievers200JSONResponse{retrieverList}, nil
}
