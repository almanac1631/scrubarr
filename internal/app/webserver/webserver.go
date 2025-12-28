package webserver

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
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

func SetupWebserver(config *koanf.Koanf, radarrRetriever *media.RadarrRetriever, delugeRetriever *torrentclients.DelugeRetriever, rtorrentRetriever *torrentclients.RtorrentRetriever) http.Handler {
	templateCache, err := NewTemplateCache()
	if err != nil {
		slog.Error("could not create template cache", "error", err)
		os.Exit(1)
	}
	handler, err := newHandler(config, templateCache, radarrRetriever, delugeRetriever, rtorrentRetriever)
	router := http.NewServeMux()
	if err != nil {
		slog.Error("could not create webserver handler", "error", err)
		os.Exit(1)
	}
	router.HandleFunc("/media", handler.handleMediaEndpoint)
	router.HandleFunc("/media/entries", handler.handleMediaEntriesEndpoint)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/media", http.StatusSeeOther)
	})
	router.Handle("/assets/", http.FileServer(http.FS(internal.Assets)))
	pathPrefix := config.String("general.path_prefix")
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
