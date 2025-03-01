package scrubarr

import (
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/internal/pkg/entrymappings/sqlite"
	"github.com/almanac1631/scrubarr/internal/pkg/retriever_bundled/simple"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
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

	entryRetrievers, err := initializeEntryRetrievers(k)
	if err != nil {
		panic(err)
	}

	dbFile := k.MustString("general.db_file")
	bundledEntryRetriever := simple.BundledEntryRetriever(k.MustStrings("general.allowed_file_endings"))
	entryMappingManager, err := sqlite.NewEntryMappingManager(entryRetrievers, bundledEntryRetriever, dbFile)
	if err != nil {
		panic(err)
	}

	listener, err := webserver.SetupListener(k)
	if err != nil {
		slog.Error("could not setup web server listener", "error", err)
		os.Exit(1)
	}
	router, err := webserver.SetupWebserver(k, entryMappingManager)
	if err != nil {
		slog.Error("could not setup web server router", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("refreshing entry mappings...")
		err := entryMappingManager.RefreshEntryMappings()
		if err != nil {
			slog.Error("could not refresh entry mappings", "error", err)
		} else {
			slog.Info("refreshed entry mappings")
		}
	}()

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
