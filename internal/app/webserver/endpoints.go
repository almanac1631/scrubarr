package webserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
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
		Name:              string(entryMapping.Name),
		RetrieverFindings: findings,
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
