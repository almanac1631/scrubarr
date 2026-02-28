package webserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"

	internal "github.com/almanac1631/scrubarr/web"
	"github.com/knadh/koanf/v2"
)

func SetupListener(config *koanf.Koanf) (net.Listener, error) {
	network := config.MustString("general.listen_network")
	addr := config.MustString("general.listen_addr")
	listener, err := net.Listen(network, addr)
	if err != nil {
		return nil, fmt.Errorf("could not listen on %s/%s: %w", network, addr, err)
	}
	slog.Info("Listening on webserver interface", "addr", addr)
	return listener, nil
}

func SetupWebserver(config *koanf.Koanf, version string, inventoryService InventoryService, quotaService QuotaService) http.Handler {
	templateCache, err := NewTemplateCache()
	if err != nil {
		slog.Error("Could not create template cache.", "error", err)
		os.Exit(1)
	}
	pathPrefix := config.String("general.path_prefix")
	if !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix += "/"
	}
	if !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = "/" + pathPrefix
	}
	realIpHeaderName := config.String("general.real_ip_header_name")
	authProvider, err := GetAuthProvider(config)
	if err != nil {
		slog.Error("Could not create auth provider.", "error", err)
		os.Exit(1)
	}
	handler, err := newHandler(config, version, pathPrefix, authProvider, templateCache, inventoryService, quotaService)
	router := http.NewServeMux()
	if err != nil {
		slog.Error("Could not create webserver handler.", "error", err)
		os.Exit(1)
	}
	router.Handle("GET /assets/", http.FileServer(http.FS(internal.Assets)))
	router.HandleFunc("/login", handler.handleLogin)
	router.HandleFunc("POST /logout", handler.handleLogout)

	authorizedRouter := http.NewServeMux()
	authorizedRouter.HandleFunc("GET /quotas/disk", handler.handleDiskQuotaEndpoint)
	authorizedRouter.HandleFunc("GET /media", handler.handleMediaEndpoint)
	authorizedRouter.HandleFunc("PUT /media", handler.handleRefreshEndpoint)
	authorizedRouter.HandleFunc("GET /media/entries", handler.handleMediaEntriesEndpoint)
	authorizedRouter.HandleFunc("GET /media/entries/{id}", handler.handleMediaSeriesEndpoint)
	authorizedRouter.HandleFunc("DELETE /media/entries/{id}", handler.handleMediaDeletionEndpoint)
	authorizedRouter.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(writer, request)
			return
		}
		http.Redirect(writer, request, path.Join(handler.pathPrefix, "/media"), http.StatusSeeOther)
	})

	router.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		redirectToLogin := func(writer http.ResponseWriter) {
			http.Redirect(writer, request, path.Join(handler.pathPrefix, "/login"), http.StatusSeeOther)
		}

		sessionCookie, err := request.Cookie(sessionCookieName)
		if errors.Is(err, http.ErrNoCookie) || sessionCookie == nil || sessionCookie.Value == "" {
			redirectToLogin(writer)
			return
		}
		token := sessionCookie.Value
		tokenOk, username, err := validateToken(handler.jwtConfig.PublicKey, token)
		if !tokenOk {
			if err != nil {
				slog.Debug("Could not validate JWT", "error", err, "token", token)
			}
			redirectToLogin(writer)
			return
		}
		logger := slog.With("remote", request.RemoteAddr).With("username", username)
		request = request.WithContext(context.WithValue(request.Context(), "logger", logger))
		authorizedRouter.ServeHTTP(writer, request)
	}))

	realIpRouter := router

	if realIpHeaderName != "" {
		slog.Info("Using real ip header.", "realIpHeaderName", realIpHeaderName)
		wrapperRouter := http.NewServeMux()
		wrapperRouter.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			realIp := request.Header.Get(realIpHeaderName)
			if realIp == "" {
				slog.Warn("Could not find value for real ip header.", "realIpHeaderName", realIpHeaderName, "remote", request.RemoteAddr)
				realIp = request.RemoteAddr
			}
			request.RemoteAddr = realIp
			router.ServeHTTP(writer, request)
		}))
		realIpRouter = wrapperRouter
	}

	if pathPrefix != "" {
		slog.Info("Applying path prefix stripping.", "path_prefix", pathPrefix)
		return http.StripPrefix(strings.TrimSuffix(pathPrefix, "/"), realIpRouter)
	}
	return realIpRouter
}

func isErrAndNoBrokenPipe(err error) bool {
	if err == nil {
		return false
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return !errors.Is(opErr.Err, syscall.EPIPE) && !errors.Is(err, io.ErrClosedPipe)
	}
	return true
}

func getRequestLogger(request *http.Request) *slog.Logger {
	logger, ok := request.Context().Value("logger").(*slog.Logger)
	if !ok {
		slog.Error("Could not get request logger.", "request", request)
		return slog.With("remote", request.RemoteAddr).With("username", "<unknown>")
	}
	return logger
}
