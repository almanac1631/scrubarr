package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/common"
	"github.com/knadh/koanf/v2"
	"log/slog"
	"net/http"
)

func StartWebserver(ctx context.Context, koanf *koanf.Koanf, retrieverRegistry common.RetrieverRegistry) error {
	// Create a new router & API
	router := http.NewServeMux()
	apiServer := &apiEndpointHandler{retrieverRegistry}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		var errorStr string
		var detail string
		var formatErr *InvalidParamFormatError
		if errors.As(err, &formatErr) {
			errorStr = "request error"
			detail = err.Error()
		} else {
			errorStr = "unknown error"
			detail = "no description provided"
			slog.Error("an unkown error occurred", "err", err)
		}
		respBody, _ := json.Marshal(ErrorResponseBody{
			Error:  errorStr,
			Detail: detail,
		})
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(respBody)
	}

	serverInterface := NewStrictHandlerWithOptions(apiServer, []StrictMiddlewareFunc{}, StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  errorHandler,
		ResponseErrorHandlerFunc: errorHandler,
	})

	HandlerWithOptions(serverInterface, StdHTTPServerOptions{
		BaseURL:    "/api",
		BaseRouter: router,
		Middlewares: []MiddlewareFunc{
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					next.ServeHTTP(w, req)
				})
			},
		},
	})

	router.Handle("/", http.FileServer(http.Dir("./web/")))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8888), router); err != nil {
		panic(err)
	}
	return nil
}
