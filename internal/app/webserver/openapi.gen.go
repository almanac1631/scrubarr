//go:build go1.22

// Package webserver provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for RetrieverCategory.
const (
	RetrieverCategoryArrApp        RetrieverCategory = "arr_app"
	RetrieverCategoryFolder        RetrieverCategory = "folder"
	RetrieverCategoryTorrentClient RetrieverCategory = "torrent_client"
)

// Defines values for RetrieverSoftwareName.
const (
	RetrieverSoftwareNameDeluge   RetrieverSoftwareName = "deluge"
	RetrieverSoftwareNameFolder   RetrieverSoftwareName = "folder"
	RetrieverSoftwareNameRadarr   RetrieverSoftwareName = "radarr"
	RetrieverSoftwareNameRtorrent RetrieverSoftwareName = "rtorrent"
	RetrieverSoftwareNameSonarr   RetrieverSoftwareName = "sonarr"
)

// Defines values for GetEntryMappingsParamsFilter.
const (
	CompleteEntries   GetEntryMappingsParamsFilter = "complete_entries"
	IncompleteEntries GetEntryMappingsParamsFilter = "incomplete_entries"
)

// Defines values for GetEntryMappingsParamsSortBy.
const (
	DateAddedAsc  GetEntryMappingsParamsSortBy = "date_added_asc"
	DateAddedDesc GetEntryMappingsParamsSortBy = "date_added_desc"
	NameAsc       GetEntryMappingsParamsSortBy = "name_asc"
	NameDesc      GetEntryMappingsParamsSortBy = "name_desc"
	SizeAsc       GetEntryMappingsParamsSortBy = "size_asc"
	SizeDesc      GetEntryMappingsParamsSortBy = "size_desc"
)

// EntryMapping defines model for EntryMapping.
type EntryMapping struct {
	// DateAdded The date and time this entry was added.
	DateAdded time.Time `json:"dateAdded"`

	// Id The unique identifier of this entry.
	Id string `json:"id"`

	// Name The name of this entry.
	Name              string                               `json:"name"`
	RetrieverFindings []EntryMappingRetrieverFindingsInner `json:"retrieverFindings"`

	// Size The size of this entry in bytes.
	Size int64 `json:"size"`
}

// EntryMappingRetrieverFindingsInner defines model for EntryMapping_retrieverFindings_inner.
type EntryMappingRetrieverFindingsInner struct {
	// Id Id used to identify retrievers.
	Id RetrieverId `json:"id"`
}

// ErrorResponseBody defines model for ErrorResponseBody.
type ErrorResponseBody struct {
	Detail string `json:"detail"`
	Error  string `json:"error"`
}

// Info defines model for Info.
type Info struct {
	// Commit The commit hash of the application.
	Commit string `json:"commit"`

	// Version The version of the application.
	Version string `json:"version"`
}

// LoginRequestBody defines model for LoginRequestBody.
type LoginRequestBody struct {
	// Password The password to login with.
	Password string `json:"password"`

	// Username The username to login with.
	Username string `json:"username"`
}

// Retriever defines model for Retriever.
type Retriever struct {
	// Category The category this retriever belongs to.
	Category RetrieverCategory `json:"category"`

	// Id Id used to identify retrievers.
	Id RetrieverId `json:"id"`

	// Name The provided name used to differentiate between multiple instances of the same software retrievers.
	Name string `json:"name"`

	// SoftwareName The name of the retriever's software.
	SoftwareName RetrieverSoftwareName `json:"softwareName"`
}

// RetrieverCategory The category this retriever belongs to.
type RetrieverCategory string

// RetrieverSoftwareName The name of the retriever's software.
type RetrieverSoftwareName string

// RetrieverId Id used to identify retrievers.
type RetrieverId = string

// Stats defines model for Stats.
type Stats struct {
	DiskSpace StatsDiskSpace `json:"diskSpace"`
}

// StatsDiskSpace defines model for Stats_diskSpace.
type StatsDiskSpace struct {
	// BytesTotal The total number of bytes available.
	BytesTotal int64 `json:"bytesTotal"`

	// BytesUsed The number of bytes used.
	BytesUsed int64 `json:"bytesUsed"`
}

// GetEntryMappings200Response defines model for getEntryMappings_200_response.
type GetEntryMappings200Response struct {
	Entries []EntryMapping `json:"entries"`

	// TotalAmount The total amount of entries that could be returned for the provided filter.
	TotalAmount int `json:"totalAmount"`
}

// GetInfo200Response defines model for getInfo_200_response.
type GetInfo200Response struct {
	Info Info `json:"info"`
}

// GetRetrievers200Response defines model for getRetrievers_200_response.
type GetRetrievers200Response struct {
	Retrievers []Retriever `json:"retrievers"`
}

// GetStats200Response defines model for getStats_200_response.
type GetStats200Response struct {
	Stats Stats `json:"stats"`
}

// Login200Response defines model for login_200_response.
type Login200Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// RefreshEntryMappings200Response defines model for refreshEntryMappings_200_response.
type RefreshEntryMappings200Response struct {
	// Message The status message to display.
	Message string `json:"message"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = ErrorResponseBody

// GetEntryMappingsParams defines parameters for GetEntryMappings.
type GetEntryMappingsParams struct {
	// Page The page number to display.
	Page int `form:"page" json:"page"`

	// PageSize The amount of items to display per each page.
	PageSize int `form:"pageSize" json:"pageSize"`

	// Filter The filter to apply before returning the entries.
	Filter *GetEntryMappingsParamsFilter `form:"filter,omitempty" json:"filter,omitempty"`

	// SortBy The criteria to sort the entries by.
	SortBy *GetEntryMappingsParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// Name The name of the entry to search for.
	Name *string `form:"name,omitempty" json:"name,omitempty"`
}

// GetEntryMappingsParamsFilter defines parameters for GetEntryMappings.
type GetEntryMappingsParamsFilter string

// GetEntryMappingsParamsSortBy defines parameters for GetEntryMappings.
type GetEntryMappingsParamsSortBy string

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequestBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get a list of entry mappings.
	// (GET /entry-mappings)
	GetEntryMappings(w http.ResponseWriter, r *http.Request, params GetEntryMappingsParams)
	// Trigger a refresh of the entry mappings.
	// (POST /entry-mappings)
	RefreshEntryMappings(w http.ResponseWriter, r *http.Request)
	// Get the information about the application.
	// (GET /info)
	GetInfo(w http.ResponseWriter, r *http.Request)
	// Login to the application using the provided credentials.
	// (POST /login)
	Login(w http.ResponseWriter, r *http.Request)
	// Get a list of retrievers.
	// (GET /retrievers)
	GetRetrievers(w http.ResponseWriter, r *http.Request)
	// Get the statistics of the application.
	// (GET /stats)
	GetStats(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetEntryMappings operation middleware
func (siw *ServerInterfaceWrapper) GetEntryMappings(w http.ResponseWriter, r *http.Request) {

	var err error

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	// Parameter object where we will unmarshal all parameters from the context
	var params GetEntryMappingsParams

	// ------------- Required query parameter "page" -------------

	if paramValue := r.URL.Query().Get("page"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "page"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "page", r.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "page", Err: err})
		return
	}

	// ------------- Required query parameter "pageSize" -------------

	if paramValue := r.URL.Query().Get("pageSize"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "pageSize"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "pageSize", r.URL.Query(), &params.PageSize)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "pageSize", Err: err})
		return
	}

	// ------------- Optional query parameter "filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "filter", r.URL.Query(), &params.Filter)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filter", Err: err})
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, false, "name", r.URL.Query(), &params.Name)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetEntryMappings(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// RefreshEntryMappings operation middleware
func (siw *ServerInterfaceWrapper) RefreshEntryMappings(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RefreshEntryMappings(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetInfo operation middleware
func (siw *ServerInterfaceWrapper) GetInfo(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInfo(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetRetrievers operation middleware
func (siw *ServerInterfaceWrapper) GetRetrievers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetRetrievers(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetStats operation middleware
func (siw *ServerInterfaceWrapper) GetStats(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStats(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/entry-mappings", wrapper.GetEntryMappings)
	m.HandleFunc("POST "+options.BaseURL+"/entry-mappings", wrapper.RefreshEntryMappings)
	m.HandleFunc("GET "+options.BaseURL+"/info", wrapper.GetInfo)
	m.HandleFunc("POST "+options.BaseURL+"/login", wrapper.Login)
	m.HandleFunc("GET "+options.BaseURL+"/retrievers", wrapper.GetRetrievers)
	m.HandleFunc("GET "+options.BaseURL+"/stats", wrapper.GetStats)

	return m
}

type ErrorResponseJSONResponse ErrorResponseBody

type GetEntryMappingsRequestObject struct {
	Params GetEntryMappingsParams
}

type GetEntryMappingsResponseObject interface {
	VisitGetEntryMappingsResponse(w http.ResponseWriter) error
}

type GetEntryMappings200JSONResponse GetEntryMappings200Response

func (response GetEntryMappings200JSONResponse) VisitGetEntryMappingsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetEntryMappings4XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetEntryMappings4XXJSONResponse) VisitGetEntryMappingsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetEntryMappings5XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetEntryMappings5XXJSONResponse) VisitGetEntryMappingsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type RefreshEntryMappingsRequestObject struct {
}

type RefreshEntryMappingsResponseObject interface {
	VisitRefreshEntryMappingsResponse(w http.ResponseWriter) error
}

type RefreshEntryMappings200JSONResponse RefreshEntryMappings200Response

func (response RefreshEntryMappings200JSONResponse) VisitRefreshEntryMappingsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type RefreshEntryMappings4XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response RefreshEntryMappings4XXJSONResponse) VisitRefreshEntryMappingsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type RefreshEntryMappings5XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response RefreshEntryMappings5XXJSONResponse) VisitRefreshEntryMappingsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetInfoRequestObject struct {
}

type GetInfoResponseObject interface {
	VisitGetInfoResponse(w http.ResponseWriter) error
}

type GetInfo200JSONResponse GetInfo200Response

func (response GetInfo200JSONResponse) VisitGetInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetInfo4XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetInfo4XXJSONResponse) VisitGetInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetInfo5XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetInfo5XXJSONResponse) VisitGetInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type LoginRequestObject struct {
	Body *LoginJSONRequestBody
}

type LoginResponseObject interface {
	VisitLoginResponse(w http.ResponseWriter) error
}

type Login200JSONResponse Login200Response

func (response Login200JSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Login401JSONResponse ErrorResponseBody

func (response Login401JSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type Login4XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response Login4XXJSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type Login5XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response Login5XXJSONResponse) VisitLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetRetrieversRequestObject struct {
}

type GetRetrieversResponseObject interface {
	VisitGetRetrieversResponse(w http.ResponseWriter) error
}

type GetRetrievers200JSONResponse GetRetrievers200Response

func (response GetRetrievers200JSONResponse) VisitGetRetrieversResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRetrievers4XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetRetrievers4XXJSONResponse) VisitGetRetrieversResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetRetrievers5XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetRetrievers5XXJSONResponse) VisitGetRetrieversResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetStatsRequestObject struct {
}

type GetStatsResponseObject interface {
	VisitGetStatsResponse(w http.ResponseWriter) error
}

type GetStats200JSONResponse GetStats200Response

func (response GetStats200JSONResponse) VisitGetStatsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetStats4XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetStats4XXJSONResponse) VisitGetStatsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetStats5XXJSONResponse struct {
	Body       ErrorResponseBody
	StatusCode int
}

func (response GetStats5XXJSONResponse) VisitGetStatsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get a list of entry mappings.
	// (GET /entry-mappings)
	GetEntryMappings(ctx context.Context, request GetEntryMappingsRequestObject) (GetEntryMappingsResponseObject, error)
	// Trigger a refresh of the entry mappings.
	// (POST /entry-mappings)
	RefreshEntryMappings(ctx context.Context, request RefreshEntryMappingsRequestObject) (RefreshEntryMappingsResponseObject, error)
	// Get the information about the application.
	// (GET /info)
	GetInfo(ctx context.Context, request GetInfoRequestObject) (GetInfoResponseObject, error)
	// Login to the application using the provided credentials.
	// (POST /login)
	Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error)
	// Get a list of retrievers.
	// (GET /retrievers)
	GetRetrievers(ctx context.Context, request GetRetrieversRequestObject) (GetRetrieversResponseObject, error)
	// Get the statistics of the application.
	// (GET /stats)
	GetStats(ctx context.Context, request GetStatsRequestObject) (GetStatsResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetEntryMappings operation middleware
func (sh *strictHandler) GetEntryMappings(w http.ResponseWriter, r *http.Request, params GetEntryMappingsParams) {
	var request GetEntryMappingsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetEntryMappings(ctx, request.(GetEntryMappingsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetEntryMappings")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetEntryMappingsResponseObject); ok {
		if err := validResponse.VisitGetEntryMappingsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// RefreshEntryMappings operation middleware
func (sh *strictHandler) RefreshEntryMappings(w http.ResponseWriter, r *http.Request) {
	var request RefreshEntryMappingsRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.RefreshEntryMappings(ctx, request.(RefreshEntryMappingsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "RefreshEntryMappings")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(RefreshEntryMappingsResponseObject); ok {
		if err := validResponse.VisitRefreshEntryMappingsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetInfo operation middleware
func (sh *strictHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	var request GetInfoRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetInfo(ctx, request.(GetInfoRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetInfo")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetInfoResponseObject); ok {
		if err := validResponse.VisitGetInfoResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Login operation middleware
func (sh *strictHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequestObject

	var body LoginJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Login(ctx, request.(LoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Login")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(LoginResponseObject); ok {
		if err := validResponse.VisitLoginResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetRetrievers operation middleware
func (sh *strictHandler) GetRetrievers(w http.ResponseWriter, r *http.Request) {
	var request GetRetrieversRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetRetrievers(ctx, request.(GetRetrieversRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetRetrievers")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetRetrieversResponseObject); ok {
		if err := validResponse.VisitGetRetrieversResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetStats operation middleware
func (sh *strictHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var request GetStatsRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetStats(ctx, request.(GetStatsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetStats")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetStatsResponseObject); ok {
		if err := validResponse.VisitGetStatsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
