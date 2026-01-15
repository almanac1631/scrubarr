package scrubarr

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/media"
	"github.com/almanac1631/scrubarr/pkg/torrentclients"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
)

var (
	version = "<no version>"
	commit  = "<no commit>"
)

var (
	k                           = koanf.New(".")
	configPath, logLevel        string
	saveCache, useCache, dryRun bool
)

var rootCmd = &cobra.Command{
	Use:     "scrubarr",
	Short:   "scrubarr is a tool to track and delete files safely on *arr instances.",
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupLogging()
		slog.Info("Scrubarr start initiated.", "version", version, "commit", commit)
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the scrubarr server",
	Run:   serve,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "level", "info", "slog level to use")
	serveCmd.Flags().StringVar(&configPath, "config", "./config.toml", "path to config file")
	serveCmd.Flags().StringVar(&logLevel, "level", "info", "log level to use")
	serveCmd.Flags().BoolVar(&saveCache, "save-cache", false, "save cache to disk")
	serveCmd.Flags().BoolVar(&useCache, "use-cache", false, "use previously saved cache for retrievers")
	serveCmd.Flags().BoolVar(&dryRun, "dry-run", false, "enable dry run mode to prevent actual file/torrent deletion")
	rootCmd.AddCommand(serveCmd)
}

func StartApp() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func serve(cmd *cobra.Command, args []string) {
	err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	if dryRun {
		slog.Info("Running in dry run mode. No files nor torrents will be deleted.")
	}

	listener, err := webserver.SetupListener(k)
	if err != nil {
		slog.Error("Could not setup web server listener.", "error", err)
		os.Exit(1)
	}

	radarrRetriever, err := media.NewRadarrRetriever(
		k.MustString("connections.radarr.hostname"),
		k.MustString("connections.radarr.api_key"),
		dryRun,
	)
	if err != nil {
		slog.Error("Could not setup radarr retriever", "error", err)
		os.Exit(1)
	}

	sonarrRetriever, err := media.NewSonarrRetriever(
		k.MustString("connections.sonarr.hostname"),
		k.MustString("connections.sonarr.api_key"),
		dryRun,
	)
	if err != nil {
		slog.Error("Could not setup sonarr retriever", "error", err)
		os.Exit(1)
	}

	mediaManager := media.NewDefaultMediaManager(radarrRetriever, sonarrRetriever)

	delugeRetriever, err := torrentclients.NewDelugeRetriever(
		k.MustString("connections.deluge.hostname"),
		uint(k.MustInt("connections.deluge.port")),
		k.MustString("connections.deluge.username"),
		k.MustString("connections.deluge.password"),
		dryRun,
	)
	if err != nil {
		slog.Error("Could not setup deluge retriever", "error", err)
		os.Exit(1)
	}

	rtorrentRetriever, err := torrentclients.NewRtorrentRetriever(
		k.MustString("connections.rtorrent.hostname"),
		k.MustString("connections.rtorrent.username"),
		k.MustString("connections.rtorrent.password"),
		dryRun,
	)
	if err != nil {
		slog.Error("Could not setup rtorrent retriever", "error", err)
		os.Exit(1)
	}

	torrentManager := torrentclients.NewDefaultTorrentManager(delugeRetriever, rtorrentRetriever)

	refreshInterval := k.Duration("general.refresh_interval")

	refreshCaches := func() {
		slog.Debug("Refreshing retriever data...")
		if err = warmupCaches(saveCache, useCache, mediaManager, torrentManager); err != nil {
			slog.Error("Could not setup retriever caches.", "error", err)
			os.Exit(1)
		}
		slog.Debug("Refreshed retriever data.")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	if refreshInterval != 0 {
		slog.Info("Starting retriever refresh automation...", "interval", refreshInterval)
		go func() {
			for {
				select {
				case <-time.After(refreshInterval):
					refreshCaches()
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	refreshCaches()
	slog.Info("Refreshed retriever caches. Setting up webserver...")

	router := webserver.SetupWebserver(k, version, mediaManager, torrentManager)

	slog.Info("Successfully set up webserver. Waiting for incoming connections...")

	go func() {
		exitChan := make(chan os.Signal, 1)
		signal.Notify(exitChan, os.Interrupt)
		<-exitChan
		slog.Info("Received exit signal. Shutting down...")
		cancelFunc()
		if err := listener.Close(); err != nil {
			slog.Error("Could not close listener.", "error", err)
		}
		slog.Info("Goodbye!")
	}()

	err = http.Serve(listener, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) && err.Error() != "" {
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Err.Error() != "use of closed network connection" {
			slog.Error("Could not serve webserver.", "error", err, "errType", fmt.Sprintf("%T", err))
		}
	}
}

func LoadConfig(configPath string) error {
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	return nil
}
