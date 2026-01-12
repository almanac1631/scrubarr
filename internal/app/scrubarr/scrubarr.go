package scrubarr

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/media"
	"github.com/almanac1631/scrubarr/pkg/torrentclients"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var (
	version = "<no version>"
	commit  = "<no commit>"
)

var k = koanf.New(".")

func StartApp() {
	f := flag.NewFlagSet("scrubarr", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println("Usage: scrubarr [options...]")
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	configPath := f.String("config", "./config.toml", "path to config file")
	logLevel := f.String("level", "info", "log level to use")
	saveCache := f.Bool("save-cache", false, "save cache for retriever responses")
	useCache := f.Bool("use-cache", false, "use previously saved cache for retrievers")
	dryRun := f.Bool("dry-run", true, "enable dry run mode to prevent actual file/torrent deletion")
	err := f.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	var level slog.Level
	err = level.UnmarshalText([]byte(*logLevel))
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	err = LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	slog.Info("starting scrubarr...", "version", version, "commit", commit)

	if *dryRun {
		slog.Info("Running in dry run mode. No files nor torrents will be deleted.")
	}

	listener, err := webserver.SetupListener(k)
	if err != nil {
		slog.Error("could not setup web server listener", "error", err)
		os.Exit(1)
	}

	radarrRetriever, err := media.NewRadarrRetriever(k.MustString("connections.radarr.hostname"), k.MustString("connections.radarr.api_key"), *dryRun)
	if err != nil {
		slog.Error("could not setup radarr retriever", "error", err)
		os.Exit(1)
	}

	sonarrRetriever, err := media.NewSonarrRetriever(k.MustString("connections.sonarr.hostname"), k.MustString("connections.sonarr.api_key"), *dryRun)
	if err != nil {
		slog.Error("could not setup sonarr retriever", "error", err)
		os.Exit(1)
	}

	mediaManager := media.NewDefaultMediaManager(radarrRetriever, sonarrRetriever)

	delugeRetriever, err := torrentclients.NewDelugeRetriever(
		k.MustString("connections.deluge.hostname"),
		uint(k.MustInt("connections.deluge.port")),
		k.MustString("connections.deluge.username"),
		k.MustString("connections.deluge.password"),
		*dryRun,
	)
	if err != nil {
		slog.Error("could not setup deluge retriever", "error", err)
		os.Exit(1)
	}

	rtorrentRetriever, err := torrentclients.NewRtorrentRetriever(
		k.MustString("connections.rtorrent.hostname"),
		k.MustString("connections.rtorrent.username"),
		k.MustString("connections.rtorrent.password"),
		*dryRun,
	)
	if err != nil {
		slog.Error("could not setup rtorrent retriever", "error", err)
		os.Exit(1)
	}

	torrentManager := torrentclients.NewDefaultTorrentManager(delugeRetriever, rtorrentRetriever)

	slog.Debug("Warming up retriever caches...")
	if err = warmupCaches(*saveCache, *useCache, mediaManager, torrentManager); err != nil {
		slog.Error("could not setup retriever caches", "error", err)
		os.Exit(1)
	}
	slog.Info("Refreshed retriever caches. Setting up webserver...")

	router := webserver.SetupWebserver(k, mediaManager, torrentManager)

	slog.Info("Successfully set up webserver. Waiting for incoming connections...")

	go func() {
		exitChan := make(chan os.Signal, 1)
		signal.Notify(exitChan, os.Interrupt)
		<-exitChan
		slog.Info("received exit signal. shutting down...")
		if err := listener.Close(); err != nil {
			slog.Error("could not close listener", "error", err)
		}
		slog.Info("bye!")
	}()

	err = http.Serve(listener, router)
	if err != nil && !errors.Is(err, http.ErrServerClosed) && err.Error() != "" {
		var opErr *net.OpError
		if errors.As(err, &opErr) && opErr.Err.Error() != "use of closed network connection" {
			slog.Error("could not serve webserver", "error", err, "errType", fmt.Sprintf("%T", err))
		}
	}
}

func LoadConfig(configPath string) error {
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	return nil
}
