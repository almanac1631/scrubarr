package scrubarr

import (
	"context"
	"errors"
	"fmt"
	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"log/slog"
	"net/http"
	"os"
)

var (
	version, commit string
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
	err = registerRetrievers(k)
	retrieverRegistry.RefreshCachedEntryMapping()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = webserver.StartWebserver(ctx, k, retrieverRegistry)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func LoadConfig(configPath string) error {
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	return nil
}

// webserver in background, catch ctrl+c
// specify config file as param
