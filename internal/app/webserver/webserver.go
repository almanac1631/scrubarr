package webserver

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path"
	"syscall"

	"github.com/almanac1631/scrubarr/pkg/media"
	"github.com/almanac1631/scrubarr/pkg/torrentclients"
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
	slog.Info("listen on webserver interface", "network", network, "addr", addr)
	return listener, nil
}

func SetupWebserver(config *koanf.Koanf, radarrRetriever *media.RadarrRetriever, sonarrRetriever *media.SonarrRetriever, delugeRetriever *torrentclients.DelugeRetriever, rtorrentRetriever *torrentclients.RtorrentRetriever) http.Handler {
	templateCache, err := NewTemplateCache()
	if err != nil {
		slog.Error("could not create template cache", "error", err)
		os.Exit(1)
	}
	pathPrefix := config.String("general.path_prefix")
	handler, err := newHandler(config, pathPrefix, templateCache, radarrRetriever, sonarrRetriever, delugeRetriever, rtorrentRetriever)
	router := http.NewServeMux()
	if err != nil {
		slog.Error("could not create webserver handler", "error", err)
		os.Exit(1)
	}
	router.Handle("/assets/", http.FileServer(http.FS(internal.Assets)))
	router.HandleFunc("/login", handler.handleLogin)
	router.HandleFunc("POST /logout", handler.handleLogout)

	authorizedRouter := http.NewServeMux()
	authorizedRouter.HandleFunc("/media", handler.handleMediaEndpoint)
	authorizedRouter.HandleFunc("/media/entries", handler.handleMediaEntriesEndpoint)
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
		tokenOk, err := validateToken(handler.jwtConfig.PublicKey, token)
		if !tokenOk {
			if err != nil {
				slog.Debug("Could not validate JWT", "error", err, "token", token)
			}
			redirectToLogin(writer)
			return
		}
		authorizedRouter.ServeHTTP(writer, request)
	}))

	if pathPrefix != "" {
		slog.Info("stripping path prefix", "path_prefix", pathPrefix)
		return http.StripPrefix(pathPrefix, router)
	}
	return router
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
