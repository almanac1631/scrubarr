package webserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"log/slog"
	"net/http"
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

	sortBy, err := parseSortBy(request.Params.SortBy)
	if err != nil {
		return nil, err
	}

	var name string
	if request.Params.Name != nil {
		name = *request.Params.Name
	}

	entryMappings, totalAmount, err := a.entryMappingManager.GetEntryMappings(page, pageSize, filter, sortBy, name)
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

func parseSortBy(by *GetEntryMappingsParamsSortBy) (common.EntryMappingSortBy, error) {
	if by == nil {
		return common.EntryMappingSortByNoSort, nil
	}
	switch *by {
	case DateAddedAsc:
		return common.EntryMappingSortByDateAsc, nil
	case DateAddedDesc:
		return common.EntryMappingSortByDateDesc, nil
	case SizeAsc:
		return common.EntryMappingSortBySizeAsc, nil
	case SizeDesc:
		return common.EntryMappingSortBySizeDesc, nil
	case NameAsc:
		return common.EntryMappingSortByNameAsc, nil
	case NameDesc:
		return common.EntryMappingSortByNameDesc, nil
	default:
		return -1, &InvalidParamFormatError{
			ParamName: "sortBy",
			Err:       fmt.Errorf("invalid sortBy parameter: %q", *by),
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

func getResponseEntryMappingFromPresencePairs(entryMapping *common.EntryMapping) EntryMapping {
	findings := make([]EntryMappingRetrieverFindingsInner, 0)
	for _, retrieverInfo := range entryMapping.RetrieversFound {
		retrieverId := RetrieverId(retrieverInfo.Id())
		finding := EntryMappingRetrieverFindingsInner{
			Id: retrieverId,
		}
		findings = append(findings, finding)
	}
	return EntryMapping{
		Id:                entryMapping.Id,
		Name:              string(entryMapping.Name),
		DateAdded:         entryMapping.DateAdded,
		Size:              entryMapping.Size,
		RetrieverFindings: findings,
	}
}

func (a ApiEndpointHandler) DeleteEntryMapping(ctx context.Context, request DeleteEntryMappingRequestObject) (DeleteEntryMappingResponseObject, error) {
	_, err := a.entryMappingManager.GetEntryMappingById(request.EntryId)
	if errors.Is(err, common.ErrEntryMappingNotFound) {
		return DeleteEntryMapping4XXJSONResponse{
			ErrorResponseBody{
				Error:  http.StatusText(http.StatusNotFound),
				Detail: fmt.Sprintf("no entry mapping with id %q found", request.EntryId),
			},
			http.StatusNotFound,
		}, nil
	} else if err != nil {
		return nil, err
	}
	err = a.entryMappingManager.DeleteEntryMappingById(request.EntryId)
	if err != nil {
		return nil, fmt.Errorf("could not delete entry mapping with id %q: %w", request.EntryId, err)
	}
	return DeleteEntryMapping200JSONResponse{
		Message: "ok",
	}, nil
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

func (a ApiEndpointHandler) GetStats(ctx context.Context, request GetStatsRequestObject) (GetStatsResponseObject, error) {
	bytesTotal, bytesUsed := int64(-1), int64(-1)
	var err error
	if a.statsRetriever != nil {
		bytesTotal, bytesUsed, err = a.statsRetriever.GetDiskStats()
		if err != nil {
			slog.Error("error getting disk stats", "error", err)
			return nil, err
		}
	}
	return GetStats200JSONResponse{Stats{
		DiskSpace: StatsDiskSpace{
			BytesTotal: bytesTotal,
			BytesUsed:  bytesUsed,
		},
	}}, nil
}

func (a ApiEndpointHandler) GetInfo(ctx context.Context, request GetInfoRequestObject) (GetInfoResponseObject, error) {
	return GetInfo200JSONResponse{
		Info{
			Commit:  a.info.Commit,
			Version: a.info.Version,
		},
	}, nil
}
