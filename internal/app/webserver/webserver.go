package webserver

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/common"
	"github.com/knadh/koanf/v2"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed all:content
var content embed.FS

func StartWebserver(ctx context.Context, koanf *koanf.Koanf, retrieverRegistry common.RetrieverRegistry) error {
	// Create a new router & API
	router := http.NewServeMux()
	apiServer, err := NewApiEndpointHandler(retrieverRegistry, koanf)
	if err != nil {
		return fmt.Errorf("could not create api endpoint handler: %w", err)
	}

	errorHandlerFunc := func(isRequest bool) func(http.ResponseWriter, *http.Request, error) {
		handler := func(w http.ResponseWriter, r *http.Request, err error) {
			var errorStr string
			var detail string
			var formatErr *InvalidParamFormatError
			if errors.As(err, &formatErr) {
				errorStr = "request error"
				detail = err.Error()
			} else {
				errorStr = "unknown error"
				detail = "no description provided"
				slog.Error("an unknown error occurred", "err", err, "errType", fmt.Sprintf("%T", err))
			}
			respBody, _ := json.Marshal(ErrorResponseBody{
				Error:  errorStr,
				Detail: detail,
			})
			header := http.StatusInternalServerError
			if isRequest {
				header = http.StatusBadRequest
			}
			w.WriteHeader(header)
			_, _ = w.Write(respBody)
		}
		return handler
	}

	serverInterface := NewStrictHandlerWithOptions(apiServer, []StrictMiddlewareFunc{}, StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  errorHandlerFunc(true),
		ResponseErrorHandlerFunc: errorHandlerFunc(false),
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
	serveFrontendFiles(router)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 8888), router); err != nil {
		panic(err)
	}
	return nil
}

func serveFrontendFiles(router *http.ServeMux) {
	fsys, err := fs.Sub(content, "dist")
	if err != nil {
		slog.Error("failed to load embedded embedded files", "err", err)
		os.Exit(1)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, "/")
		if filePath == "" {
			filePath = "index.html"
		}
		data, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		mimeType := mime.TypeByExtension(path.Ext(filePath))
		if mimeType == "" {
			mimeType = http.DetectContentType(data)
		}
		w.Header().Set("Content-Type", mimeType)

		_, err = w.Write(data)
		if err != nil {
			slog.Error("failed to write response", "err", err)
		}
	})
}
