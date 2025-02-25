package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/arr_apps"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/folder_scanning"
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval/torrent_clients"
	"log/slog"
)

func (a ApiEndpointHandler) GetEntryMappings(ctx context.Context, request GetEntryMappingsRequestObject) (GetEntryMappingsResponseObject, error) {
	page := request.Params.Page
	pageSize := request.Params.PageSize
	if err := validatePage(request.Params.Page); err != nil {
		return nil, err
	}
	if err := validatePageSize(request.Params.PageSize); err != nil {
		return nil, err
	}
	filter, err := parseRequestFilter(request.Params.Filter)
	if err != nil {
		return nil, err
	}

	entryMappings, totalAmount, err := a.entryMappingManager.GetEntryMappings(page, pageSize, filter)
	if err != nil {
		return nil, err
	}
	var entryMappingRespList []EntryMapping
	for _, entryMapping := range entryMappings {
		entryMappingRespList = append(entryMappingRespList, getResponseEntryMappingFromPresencePairs(entryMapping))
	}

	return &GetEntryMappings200JSONResponse{
		entryMappingRespList,
		totalAmount,
	}, nil
}

func parseRequestFilter(paramFilter *GetEntryMappingsParamsFilter) (common.EntryMappingFilter, error) {
	if paramFilter == nil {
		return common.EntryMappingFilterNoFilter, nil
	}
	switch *paramFilter {
	case "incomplete_entries":
		return common.EntryMappingFilterIncompleteEntry, nil
	case "complete_entries":
		return common.EntryMappingFilterCompleteEntry, nil
	default:
		return -1, &InvalidParamFormatError{
			ParamName: "filter",
			Err:       fmt.Errorf("invalid filter parameter: %q", *paramFilter),
		}
	}
}

func validatePage(page int) error {
	if page < 1 {
		return &InvalidParamFormatError{
			ParamName: "page",
			Err:       errors.New("page param has to be greater than 1"),
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

func getResponseEntryMappingFromPresencePairs(entryMapping common.EntryMapping) EntryMapping {
	findings := make([]EntryMappingRetrieverFindingsInner, 0)
	for retrieverInfo, entry := range entryMapping.Pairs {
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
	return EntryMapping{
		Name:              string(entryMapping.Name),
		RetrieverFindings: findings,
	}
}

func getFindingValueFromEntry(entry *common.Entry) (any, error) {
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

func (a ApiEndpointHandler) RefreshEntryMappings(ctx context.Context, request RefreshEntryMappingsRequestObject) (RefreshEntryMappingsResponseObject, error) {
	err := a.entryMappingManager.RefreshEntryMappings()
	if err != nil {
		return nil, err
	}
	return &RefreshEntryMappings200JSONResponse{Message: "ok"}, nil
}

func (a ApiEndpointHandler) GetRetrievers(ctx context.Context, request GetRetrieversRequestObject) (GetRetrieversResponseObject, error) {
	var retrieverList []Retriever
	retrievers, err := a.entryMappingManager.GetRetrievers()
	if err != nil {
		return nil, err
	}
	for _, retrieverInfo := range retrievers {
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
